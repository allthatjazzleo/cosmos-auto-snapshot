package internal

import (
	"context"
	"io"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureUploader struct {
	container string
	client    *azblob.Client
}

func NewAzureUploader(azureAccountURL, azureContainer string) (*AzureUploader, error) {
	// Create a new service client with token credential
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := azblob.NewClient(azureAccountURL, credential, nil)
	if err != nil {
		return nil, err
	}

	return &AzureUploader{
		container: azureContainer,
		client:    client,
	}, nil
}

func (u *AzureUploader) Upload(reader io.Reader, filename string) error {
	ctx := context.Background()
	progressReader := &ProgressReader{
		reader: reader,
	}

	// Upload the file to the specified container with the specified blob name
	_, err := u.client.UploadStream(ctx, u.container, filename, progressReader, &azblob.UploadStreamOptions{
		BlockSize:   512 * 1024 * 1024, // 512 MiB
		Concurrency: 4,
	})
	if err != nil {
		return err
	}
	log.Println("Uploaded", filename, "to Azure storage blob", u.container)
	return nil
}
