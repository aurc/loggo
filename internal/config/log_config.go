package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"

	"gopkg.in/yaml.v3"
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
	ColorWhen []ColorWhen `json:"color-when,omitempty" yaml:"color-when,omitempty"`
}

func (k *Key) ExtractValue(m map[string]interface{}) string {
	go func() {
		r := recover()
		if r != nil {
			log.Fatalf(`failed to process %v with val %v`, r, m)
		}
	}()
	kList := strings.Split(k.Name, "/")
	var val string
	level := m
	for i, levelKey := range kList {
		lv := level[levelKey]
		if i == len(kList)-1 {
			return fmt.Sprintf("%v", lv)
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
		yamlBytes = []byte(defaultConfig)
	}
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type Type string

func (t Type) GetColor() tcell.Color {
	switch t {
	case TypeString:
		return tcell.ColorWhite
	case TypeNumber:
		return tcell.ColorBlue
	case TypeBool:
		return tcell.ColorOrange
	case TypeDateTime:
		return tcell.ColorPurple
	}
	return tcell.ColorWhite
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
    color:
      foreground: white
      background: black`
