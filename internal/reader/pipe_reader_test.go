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
		reader := MakeReader("")
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
		streamReceiver := make(chan string, 1)
		var lines []string
		reader.StreamInto(streamReceiver)
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
