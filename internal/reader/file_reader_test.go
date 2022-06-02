package reader

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFileStream_StreamInto(t *testing.T) {
	t.Run("Test Successful Stream and closure of file", func(t *testing.T) {
		tmpDir := os.TempDir()
		fileName := uuid.New().String() + ".txt"
		filePath := path.Join(tmpDir, fileName)
		file, err := os.Create(filePath)
		assert.NoError(t, err)
		assert.FileExists(t, filePath)
		assert.NoError(t, file.Close())
		fmt.Println("created temp file ", filePath)

		// Routine to write file lines
		before := time.Now().UnixMilli()
		reader := MakeReader(filePath)
		go func() {
			for i := 0; i < 10; i++ {
				file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					assert.NoError(t, err)
				}
				_, err = file.WriteString(fmt.Sprintf("line %d\n", i+1))
				assert.NoError(t, err)
				assert.NoError(t, file.Sync())
				assert.NoError(t, file.Close())
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
