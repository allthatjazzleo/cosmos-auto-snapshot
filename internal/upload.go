package internal

import "io"

type Uploader interface {
	Upload(reader io.Reader, filename string) error
}
