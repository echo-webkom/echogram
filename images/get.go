package images

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func HandleGetImageByUserId(c *fiber.Ctx) error {
	userId := c.Query("userId")
	if userId == "" {
		return c.Status(400).SendString("Add ?userId=<userId> to the URL to get an image")
	}

	accountName, accountKey, containerName := getCredentials()
	if accountName == "" || accountKey == "" || containerName == "" {
		return c.Status(500).SendString("Storage account information is not configured")
	}

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	creds, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return c.Status(500).SendString("Failed to create credential")
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	ctx := context.Background()
	listBlobs, err := containerURL.ListBlobsFlatSegment(ctx, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})
	if err != nil {
		log.Println(err)
		return c.Status(500).SendString("Failed to list blobs: " + err.Error())
	}

	var foundBlob *azblob.BlobItemInternal
	for _, blob := range listBlobs.Segment.BlobItems {
		if strings.HasPrefix(blob.Name, userId+".") && validImageType(blob.Name) {
			foundBlob = &blob
			break
		}
	}

	if foundBlob == nil {
		return c.Status(404).SendString("Image not found")
	}

	blobURL := containerURL.NewBlobURL(foundBlob.Name)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Println(err)
		return c.Status(404).SendString("Failed to download image: " + err.Error())
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 3})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
	if err != nil {
		return c.Status(500).SendString("Error reading blob data")
	}

	c.Type(filepath.Ext(foundBlob.Name))
	return c.Status(200).Send(data)
}
