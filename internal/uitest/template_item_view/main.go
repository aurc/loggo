package main

import "github.com/aurc/loggo/internal/loggo"

func main() {
	app := loggo.NewApp("")
	view := loggo.NewTemplateItemView(app, nil, nil, nil)
	app.Run(view)
}
