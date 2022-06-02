package reader

type reader struct {
	strChan chan<- string
}

// MakeReader builds a continues file/pipe streamer used to feed the logger. If
// fileName is not provided, it will attempt to consume the input from the stdin.
func MakeReader(fileName string) Reader {
	if len(fileName) > 0 {
		return &fileStream{
			fileName: fileName,
		}
	}
	return &readPipeStream{}
}

type Reader interface {
	// StreamInto feeds the strChan channel for every streamed line.
	StreamInto(strChan chan<- string) error
	// Close finalises and invalidates this stream reader.
	Close()
}
