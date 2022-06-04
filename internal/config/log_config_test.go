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
			wants:     defConfig,
		},
		{
			name:      "Valid value supplied",
			givenFile: "../../config-sample/gcp.yaml",
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
			givenFile:  "../../testdata/test3.txt",
			wants:      defConfig,
			wantsError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
						Foreground: "white",
						Background: "black",
					},
				},
				{
					MatchValue: "WARN",
					Color: Color{
						Foreground: "red",
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
		},
		{
			Name: "resource/labels/container_name",
			Type: TypeString,
		},
		{
			Name: "trace",
			Type: TypeString,
		},
		{
			Name: "jsonPayload/message",
			Type: TypeString,
		},
	},
}
