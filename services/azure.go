package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

var (
	ErrAzureBlobCorruptedData  = errors.New("azure blob corrupted data")
	ErrAzureBlobNotFound       = errors.New("azure blob not found")
	ErrAzureBlobFailedToUpload = errors.New("azure blob failed to upload")
)

type AzureBlobManager struct {
	credentials AzureCredentials
	url         *url.URL
}

func NewAzureBlobManager() (*AzureBlobManager, error) {
	credenditals := getAzureCredentials()

	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", credenditals.AccountName, credenditals.ContainerName))
	if err != nil {
		return nil, err
	}

	return &AzureBlobManager{
		credentials: credenditals,
		url:         URL,
	}, nil
}

func (s *AzureBlobManager) Get(filename string) ([]byte, error) {
	blobURL, err := s.getBlobURL(filename)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, ErrAzureBlobNotFound
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 3})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
	if err != nil {
		return nil, ErrAzureBlobCorruptedData
	}

	return data, nil
}

func (s *AzureBlobManager) Add(filename string, file []byte) error {
	blobURL, err := s.getBlobURL(filename)
	if err != nil {
		return err
	}

	blockBlobURL := blobURL.ToBlockBlobURL()

	ctx := context.Background()

	_, err = azblob.UploadBufferToBlockBlob(ctx, file, blockBlobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return ErrAzureBlobFailedToUpload
	}

	return nil
}

func (s *AzureBlobManager) getContainerURL() (azblob.ContainerURL, error) {
	creds, err := azblob.NewSharedKeyCredential(s.credentials.AccountName, s.credentials.AccountKey)
	if err != nil {
		return azblob.ContainerURL{}, err
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*s.url, pipeline)

	return containerURL, nil
}

func (s *AzureBlobManager) getBlobURL(filename string) (azblob.BlobURL, error) {
	containerURL, err := s.getContainerURL()
	if err != nil {
		return azblob.BlobURL{}, err
	}

	blobURL := containerURL.NewBlobURL(filename)

	return blobURL, nil
}
