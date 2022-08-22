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

package main

import (
	"os"
	"path"

	"github.com/google/uuid"

	"github.com/aurc/loggo/internal/reader"
	"github.com/aurc/loggo/internal/server"
	"github.com/aurc/loggo/internal/uitest/helper"
)

func main() {
	tmpDir := os.TempDir()
	fileName := uuid.New().String() + ".txt"
	filePath := path.Join(tmpDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(filePath)
	}()
	_ = file.Close()

	inputChan := make(chan string, 1)
	rd := reader.MakeReader(filePath, inputChan)
	go func() {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		helper.JsonGenerator(file)
	}()

	if err := server.Run(&server.Settings{
		Reader: rd,
		Config: nil,
		Port:   8080,
	}); err != nil {
		panic(err)
	}
}
