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

package loggo

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aurc/loggo/internal/util"
)

const (
	parentPath = ".loggo"
	logsPath   = "logs"
)

var LogFile string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	now := time.Now().Local().Format("2006.01.02T15.04.05")
	file := fmt.Sprintf("%s.log", now)
	paramsDir := path.Join(home, parentPath, logsPath)
	if err := os.MkdirAll(paramsDir, os.ModePerm); err != nil {
		panic(err)
	}
	LogFile = path.Join(paramsDir, file)
	util.InitializeLogging(LogFile)
}
