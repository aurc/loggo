package main

import (
	"github.com/aurc/loggo/internal/config"
	"github.com/aurc/loggo/internal/loggo"
)

func main() {
	app := loggo.NewApp("")
	view := loggo.NewTemplateItemView(app, &config.Key{
		Type: config.TypeString,
		ColorWhen: []config.ColorWhen{
			{
				MatchValue: "Some String",
				Color: config.Color{
					Foreground: "white",
					Background: "purple",
				},
			},
			{
				MatchValue: "Some String",
				Color: config.Color{
					Foreground: "white",
					Background: "red",
				},
			},
		},
	}, nil, nil)
	app.Run(view)
}
