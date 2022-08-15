/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package config

import (
	"fmt"
	"sort"
	"strings"
)

func MakeConfigFromSample(sample []map[string]interface{}, mergeWith ...Key) (*Config, map[string]*Key) {
	keyMap := make(map[string]*Key)
	for i := range mergeWith {
		v := mergeWith[i]
		if _, ok := keyMap[v.Name]; !ok {
			keyMap[v.Name] = &v
		}
	}
	for _, m := range sample {
		for _, k := range extractKeys2ndDepth(m) {
			if _, ok := keyMap[k]; ok {
				continue
			}
			if k == ParseErr {
				continue
			}
			if timestamp.Contains(k) {
				keyMap[k] = timestamp.keyConfig(k)
				continue
			} else if logType.Contains(k) {
				keyMap[k] = logType.keyConfig(k)
				continue
			} else if traceId.Contains(k) {
				keyMap[k] = traceId.keyConfig(k)
				continue
			} else if message.Contains(k) {
				keyMap[k] = message.keyConfig(k)
				continue
			} else if errorKey.Contains(k) {
				keyMap[k] = errorKey.keyConfig(k)
				continue
			}
			//if _, ok := v.(map[string]interface{}); ok {
			//	continue
			//}
			keyMap[k] = &Key{
				Name: k,
				Type: TypeString,
				Color: Color{
					Foreground: "white",
					Background: "black",
				},
				MaxWidth: 25,
			}
			continue
		}
	}
	c := &Config{
		Keys: []Key{},
	}
	var orderedKeys []string
	orderedKeys = append(orderedKeys, timestamp.Keys()...)
	orderedKeys = append(orderedKeys, logType.Keys()...)
	orderedKeys = append(orderedKeys, traceId.Keys()...)
	orderedKeys = append(orderedKeys, message.Keys()...)
	orderedKeys = append(orderedKeys, errorKey.Keys()...)
	for _, v := range orderedKeys {
		if v, ok := keyMap[v]; ok {
			c.Keys = append(c.Keys, *v)
		}
	}

	var sk []string
	for k := range keyMap {
		if !timestamp.Contains(k) && !message.Contains(k) && !traceId.Contains(k) && !logType.Contains(k) && !errorKey.Contains(k) {
			sk = append(sk, k)
		}
	}
	sort.Strings(sk)
	for _, v := range sk {
		c.Keys = append(c.Keys, *keyMap[v])
	}
	return c, keyMap
}

type preBakedRule struct {
	keyMatchesAny map[string]bool
	keyConfig     func(keyName string) *Key
}

func (p preBakedRule) Contains(key string) bool {
	if _, ok := p.keyMatchesAny[key]; ok {
		return ok
	}
	return false
}

func (p preBakedRule) Keys() []string {
	var arr []string
	for k := range p.keyMatchesAny {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	return arr
}

func extractKeys2ndDepth(m map[string]interface{}) []string {
	keys := make([]string, 0)
	for k, v := range m {
		if strings.Contains(k, "/") {
			continue
		}
		if vk, ok := v.(map[string]interface{}); ok &&
			k != "http_request" &&
			k != "labels" {
			for k2 := range vk {
				if strings.Contains(k2, "/") {
					continue
				}
				keys = append(keys, fmt.Sprintf(`%s/%s`, k, k2))
			}
		} else {
			keys = append(keys, k)
		}
	}
	return keys
}

var (
	timestamp = preBakedRule{
		keyMatchesAny: map[string]bool{"timestamp": true, "time": true},
		keyConfig: func(keyName string) *Key {
			return &Key{
				Name: keyName,
				Type: TypeDateTime,
				Color: Color{
					Foreground: "purple",
					Background: "black",
				},
			}
		},
	}
	traceId = preBakedRule{
		keyMatchesAny: map[string]bool{"traceId": true},
		keyConfig: func(keyName string) *Key {
			return &Key{
				Name:     keyName,
				Type:     TypeDateTime,
				MaxWidth: 32,
				Color: Color{
					Foreground: "olive",
					Background: "black",
				},
			}
		},
	}
	logType = preBakedRule{
		keyMatchesAny: map[string]bool{"level": true, "severity": true},
		keyConfig: func(keyName string) *Key {
			return &Key{
				Name: keyName,
				Type: TypeString,
				Color: Color{
					Foreground: "white",
					Background: "black",
				},
				ColorWhen: []ColorWhen{
					{
						MatchValue: "(?i)error",
						Color: Color{
							Foreground: "red",
							Background: "black",
						},
					},
					{
						MatchValue: "(?i)info",
						Color: Color{
							Foreground: "green",
							Background: "black",
						},
					},
					{
						MatchValue: "(?i)warn",
						Color: Color{
							Foreground: "orange",
							Background: "black",
						},
					},
					{
						MatchValue: "(?i)debug",
						Color: Color{
							Foreground: "blue",
							Background: "black",
						},
					},
				},
			}
		},
	}
	message = preBakedRule{
		keyMatchesAny: map[string]bool{
			"message":             true,
			"jsonPayload/message": true,
			"http_request":        true,
		},
		keyConfig: func(keyName string) *Key {
			return &Key{
				Name:     keyName,
				Type:     TypeString,
				MaxWidth: 60,
				Color: Color{
					Foreground: "wheat",
					Background: "black",
				},
			}
		},
	}
	errorKey = preBakedRule{
		keyMatchesAny: map[string]bool{"error": true},
		keyConfig: func(keyName string) *Key {
			return &Key{
				Name:     keyName,
				Type:     TypeString,
				MaxWidth: 30,
				Color: Color{
					Foreground: "red",
					Background: "black",
				},
			}
		},
	}
)
