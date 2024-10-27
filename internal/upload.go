package internal

import (
	"io"
	"log"
)

type Uploader interface {
	Upload(reader io.Reader, filename string) error
}

type ProgressReader struct {
	reader    io.Reader
	readBytes int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.readBytes += int64(n)
	if n > 0 {
		log.Printf("Total uploaded so far: %.2f MB\n", float64(pr.readBytes)/(1024*1024))
	}
	return n, err
}
