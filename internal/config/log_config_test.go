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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeConfig(t *testing.T) {
	tests := []struct {
		name       string
		givenFile  string
		wants      Config
		wantsError bool
	}{
		{
			name:      "No file supplied, load GCP default",
			givenFile: "",
			wants:     Config{},
		},
		{
			name:      "Valid value supplied",
			givenFile: "../config-sample/gcp.yaml",
			wants:     defConfig,
		},
		{
			name:       "Non existing file",
			givenFile:  "foo",
			wants:      defConfig,
			wantsError: true,
		},
		{
			name:       "Bad format file",
			givenFile:  "../testdata/test3.txt",
			wants:      defConfig,
			wantsError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.wants.LastSavedName = test.givenFile
			c, err := MakeConfig(test.givenFile)
			if test.wantsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.wants, *c)
			}
		})
	}
}

func TestKey_ExtractValue(t *testing.T) {
	tests := []struct {
		name      string
		givenKey  *Key
		givenJson []byte
		wantValue string
	}{
		{
			name: "One level key",
			givenKey: &Key{
				Name: "value",
			},
			givenJson: []byte(`{"value":"foo"}`),
			wantValue: "foo",
		},
		{
			name: "Multi level string key",
			givenKey: &Key{
				Name: "a/b/value",
			},
			givenJson: []byte(`{"a":{"b":{"value": "foo"}}}`),
			wantValue: "foo",
		},
		{
			name: "Multi level int key",
			givenKey: &Key{
				Name: "a/b/value",
			},
			givenJson: []byte(`{"a":{"b":{"value": 1}}}`),
			wantValue: "1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := make(map[string]interface{})
			err := json.Unmarshal(test.givenJson, &m)
			assert.NoError(t, err)
			val := test.givenKey.ExtractValue(m)
			assert.Equal(t, test.wantValue, val)
		})
	}
}

var defConfig = Config{
	Keys: []Key{
		{
			Name:   "timestamp",
			Type:   TypeDateTime,
			Layout: "2006-01-02T15:04:05-0700",
			Color: Color{
				Foreground: "purple",
				Background: "black",
			},
		},
		{
			Name: "severity",
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
						Foreground: "yellow",
						Background: "black",
					},
				},
				{
					MatchValue: "DEBUG",
					Color: Color{
						Foreground: "blue",
						Background: "black",
					},
				},
			},
		},
		{
			Name: "resource/labels/container_name",
			Type: TypeString,
			Color: Color{
				Foreground: "darkgreen",
				Background: "black",
			},
		},
		{
			Name: "trace",
			Type: TypeString,
			Color: Color{
				Foreground: "white",
				Background: "black",
			},
		},
		{
			Name: "jsonPayload/message",
			Type: TypeString,
			Color: Color{
				Foreground: "white",
				Background: "black",
			},
		},
	},
}
