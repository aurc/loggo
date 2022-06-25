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
	"sort"
	"strings"
)

func MakeConfigFromSample(sample []map[string]interface{}) *Config {
	const cTimestamp = "timestamp"
	const cSeverity = "severity"
	keyMap := make(map[string]*Key)
	for _, m := range sample {
		for k, v := range m {
			if _, ok := keyMap[k]; ok {
				continue
			}
			if strings.Index(k, "/") != -1 || k == ParseErr || k == TextPayload {
				continue
			}
			if k == cTimestamp {
				keyMap[k] = &Key{
					Name: k,
					Type: TypeDateTime,
					Color: Color{
						Foreground: "purple",
						Background: "black",
					},
				}
				continue
			} else if k == cSeverity {
				keyMap[k] = &Key{
					Name: k,
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
				continue
			}
			if _, ok := v.(map[string]interface{}); ok {
				continue
			}
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
	if v, ok := keyMap[cTimestamp]; ok {
		c.Keys = append(c.Keys, *v)
	}
	if v, ok := keyMap[cSeverity]; ok {
		c.Keys = append(c.Keys, *v)
	}
	var sk []string
	for k := range keyMap {
		if k != cTimestamp && k != cSeverity {
			sk = append(sk, k)
		}
	}
	sort.Strings(sk)
	for _, v := range sk {
		c.Keys = append(c.Keys, *keyMap[v])
	}

	return c
}
