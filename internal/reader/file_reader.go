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

func (s *fileStream) StreamInto(strChan chan<- string) error {
	s.strChan = strChan
	var err error
	s.tail, err = tail.TailFile(s.fileName, tail.Config{Follow: true})
	if err != nil {
		return err
	}

	go func() {
		for line := range s.tail.Lines {
			strChan <- line.Text
		}
	}()
	return nil
}

func (s *fileStream) Close() {
	s.tail.Kill(fmt.Errorf("stopped by Close method"))
	close(s.strChan)
}
