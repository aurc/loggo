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
			if strings.Index(k, "/") != -1 {
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
							MatchValue: "ERROR",
							Color: Color{
								Foreground: "white",
								Background: "red",
							},
						},
						{
							MatchValue: "INFO",
							Color: Color{
								Foreground: "green",
								Background: "black",
							},
						},
						{
							MatchValue: "WARN",
							Color: Color{
								Foreground: "white",
								Background: "yellow",
							},
						},
						{
							MatchValue: "DEBUG",
							Color: Color{
								Foreground: "white",
								Background: "blue",
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
