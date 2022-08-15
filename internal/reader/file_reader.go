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

	"github.com/nxadm/tail"
)

type fileStream struct {
	reader
	fileName string
	tail     *tail.Tail
}

func (s *fileStream) StreamInto() error {
	var err error
	s.tail, err = tail.TailFile(s.fileName, tail.Config{Follow: true, Poll: true})
	if err != nil {
		return err
	}

	go func() {
		for line := range s.tail.Lines {
			s.strChan <- line.Text
		}
	}()
	return nil
}

func (s *fileStream) Close() {
	s.tail.Kill(fmt.Errorf("stopped by Close method"))
	close(s.strChan)
}
