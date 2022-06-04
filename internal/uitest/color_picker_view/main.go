package main

import "github.com/aurc/loggo/internal/loggo"

func main() {
	app := loggo.NewApp("")
	view := loggo.NewColourPickerView(app, "Select Color",
		func(c string) {
		}, func() {
			app.Stop()
		}, func() {
			app.Stop()
		})
	app.Run(view)
}
