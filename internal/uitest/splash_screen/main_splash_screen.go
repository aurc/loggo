package main

import "github.com/aurc/loggo/internal/loggo"

func main() {
	app := loggo.NewApp("")
	view := loggo.NewSplashScreen(app)
	app.Run(view)
}
