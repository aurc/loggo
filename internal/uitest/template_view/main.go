package main

import "github.com/aurc/loggo/internal/loggo"

func main() {
	app := loggo.NewApp("")
	view := loggo.NewTemplateView(app, nil, nil)
	app.Run(view)
}
