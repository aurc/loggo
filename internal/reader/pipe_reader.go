package reader

import (
	"bufio"
	"fmt"
	"os"
)

type readPipeStream struct {
	reader
	stop bool
}

func (s *readPipeStream) StreamInto(strChan chan<- string) error {
	s.strChan = strChan
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return fmt.Errorf("nothing in input stream")
	}

	reader := bufio.NewReader(os.Stdin)

	go func() {
		for !s.stop {
			str, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			strChan <- str
		}
	}()
	return nil
}
func (s *readPipeStream) Close() {
	s.stop = true
	close(s.strChan)
}
