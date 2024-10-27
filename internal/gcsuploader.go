package internal

import (
	"context"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

type GCSUploader struct {
	bucket string
	client *storage.Client
}

func NewGCSUploader(bucket string) (*GCSUploader, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSUploader{
		bucket: bucket,
		client: client,
	}, nil
}

func (u *GCSUploader) Upload(reader io.Reader, filename string) error {
	ctx := context.Background()
	defer u.client.Close()

	progressReader := &ProgressReader{
		reader: reader,
	}

	wc := u.client.Bucket(u.bucket).Object(filename).NewWriter(ctx)
	wc.ChunkSize = 512 * 1024 * 1024 // 512 MiB
	defer wc.Close()

	if _, err := io.Copy(wc, progressReader); err != nil {
		return err
	}

	log.Println("Uploaded", filename, "to GCS bucket", u.bucket)
	return nil
}
