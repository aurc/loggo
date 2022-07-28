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

package reader

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadPipeStream_StreamInto(t *testing.T) {
	t.Run("Test Successful Stream and closure of stdin", func(t *testing.T) {
		oldStdIn := os.Stdin
		defer func() {
			os.Stdin = oldStdIn
		}()

		// Routine to write file lines
		before := time.Now().UnixMilli()
		streamReceiver := make(chan string, 1)
		reader := MakeReader("", streamReceiver)
		r, w, err := os.Pipe()
		os.Stdin = r
		assert.NoError(t, err)
		go func() {
			for i := 0; i < 10; i++ {
				_, err := w.WriteString(fmt.Sprintf("line %d\n", i+1))
				assert.NoError(t, err)
				time.Sleep(100 * time.Millisecond)
			}
			reader.Close()
		}()
		var lines []string
		_ = reader.StreamInto()
		for {
			line, ok := <-streamReceiver
			if !ok {
				break
			}
			if len(line) > 0 {
				lines = append(lines, line)
			}
		}
		assert.Len(t, lines, 10)
		now := time.Now().UnixMilli()
		diff := (now - before) / int64(1000)
		assert.True(t, diff >= int64(1))
	})
}
