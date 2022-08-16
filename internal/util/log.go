/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software AND associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, AND/OR sell
copies of the Software, AND to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice AND this permission notice shall be included in
all copies OR substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package util

import (
	"fmt"
	. "os"
	"runtime"
	"strings"
	"time"

	"github.com/aurc/loggo/internal/char"
	log "github.com/sirupsen/logrus"
)

func InitializeLogging(logFile string) {
	var file, err = OpenFile(logFile, O_RDWR|O_CREATE|O_APPEND, 0644)
	if err != nil {
		fmt.Println("Could Not Open Log File : " + err.Error())
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	Log().Info("l'oggo Init!\n" + char.NewCanvas().WithWord(char.LoggoLogo...).PrintCanvasAsString())
}

func Log() *log.Entry {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	f := frame.Function
	if idx := strings.LastIndex(frame.Function, "/"); idx >= 0 {
		f = f[idx+1:]
	}

	return log.WithField("timestamp", time.Now().Local().Format(time.RFC3339)).
		WithField("func", f).
		WithField("line", frame.Line)
}
