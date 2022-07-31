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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"

	"gopkg.in/yaml.v3"
)

const (
	ParseErr    = "$_parseErr"
	TextPayload = "message"
)

type Config struct {
	Keys          []Key  `json:"keys" yaml:"keys"`
	LastSavedName string `json:"-" yaml:"-"`
}

func (c *Config) Save(fileName string) error {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	c.LastSavedName = fileName
	return nil
}

func (c *Config) KeyMap() map[string]*Key {
	nk := make(map[string]*Key)
	for _, k := range c.Keys {
		kp := &k
		nk[k.Name] = kp
	}
	return nk
}

type Color struct {
	Foreground string `json:"foreground" yaml:"foreground"`
	Background string `json:"background" yaml:"background"`
}

func (c *Color) GetBackgroundColor() tcell.Color {
	if len(c.Background) > 0 {
		return tcell.GetColor(strings.ToLower(c.Background))
	}
	return tcell.ColorBlack
}

func (c *Color) GetForegroundColor() tcell.Color {
	if len(c.Foreground) > 0 {
		return tcell.GetColor(strings.ToLower(c.Foreground))
	}
	return tcell.ColorWhite
}

func (c *Color) SetTextTagColor(text string) string {
	return fmt.Sprintf(`[%s:%s:]%s[-:-:]`,
		c.Foreground, c.Background, text)
}

type Match struct {
	Value string `json:"value" yaml:"value,omitempty"`
	Color Color  `json:"color" yaml:"color,omitempty"`
}

type ColorWhen struct {
	MatchValue string `json:"match-value" yaml:"match-value,omitempty"`
	Color      Color  `json:"color" yaml:"color,omitempty"`
}

type Key struct {
	Name      string      `json:"name" yaml:"name"`
	Type      Type        `json:"type" yaml:"type"`
	Layout    string      `json:"layout,omitempty" yaml:"layout,omitempty"`
	Color     Color       `json:"color,omitempty" yaml:"color,omitempty"`
	MaxWidth  int         `json:"max-width,omitempty" yaml:"max-width"`
	ColorWhen []ColorWhen `json:"color-when,omitempty" yaml:"color-when,omitempty"`
}

func GetForegroundColorName(colorable func() *Color, colorIfNone string) string {
	k := colorable()
	if k == nil || len(k.Foreground) < 0 {
		return colorIfNone
	}
	return k.Foreground
}

func GetBackgroundColorName(colorable func() *Color, colorIfNone string) string {
	k := colorable()
	if k == nil || len(k.Background) < 0 {
		return colorIfNone
	}
	return k.Background
}

func (k *Key) ExtractValue(m map[string]interface{}) string {
	kList := strings.Split(k.Name, "/")
	var val string
	level := m
	for i, levelKey := range kList {
		lv := level[levelKey]
		if lv == nil {
			return val
		}
		if i == len(kList)-1 {
			if v, ok := lv.(map[string]interface{}); ok {
				b, err := json.Marshal(v)
				if err == nil {
					return string(b)
				}
			}
			return fmt.Sprintf("%+v", lv)
		}
		level = lv.(map[string]interface{})
	}
	return val
}

func MakeConfig(file string) (*Config, error) {
	var yamlBytes []byte
	config := Config{}
	if len(file) > 0 {
		var err error
		yamlBytes, err = os.ReadFile(file)
		if err != nil {
			return nil, err
		}
	} else {
		yamlBytes = []byte("")
	}
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, err
	}
	config.LastSavedName = file
	return &config, nil
}

type Type string

func (t Type) GetColorName() string {
	switch t {
	case TypeString:
		return "white"
	case TypeNumber:
		return "blue"
	case TypeBool:
		return "orange"
	case TypeDateTime:
		return "purple"
	}
	return "lightgray"
}

func (t Type) GetColor() tcell.Color {
	return tcell.GetColor(t.GetColorName())
}

const (
	TypeString   = "string"
	TypeBool     = "bool"
	TypeNumber   = "number"
	TypeDateTime = "datetime"
)

const defaultConfig = `keys:
  - name: timestamp
    type: datetime
    layout: 2006-01-02T15:04:05-0700
    color:
      foreground: purple
      background: black
  - name: severity
    type: string
    color:
      foreground: white
      background: black
    color-when:
      - match-value: ERROR
        color:
          foreground: white
          background: red
      - match-value: INFO
        color:
          foreground: green
          background: black
      - match-value: WARN
        color:
          foreground: yellow
          background: black
      - match-value: DEBUG
        color:
          foreground: blue
          background: black
  - name: resource/labels/container_name
    type: string
    color:
      foreground: darkgreen
      background: black
  - name: trace
    type: string
    color:
      foreground: white
      background: black
  - name: jsonPayload/message
    type: string
    max-width: 40
    color:
      foreground: white
      background: black`
