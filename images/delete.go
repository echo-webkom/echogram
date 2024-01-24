package images

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func HandleDeleteImageByUserId(c *fiber.Ctx) error {
	userId := c.Query("userId")
	if userId == "" {
		return c.Status(400).SendString("Add ?userId=<userId> to the URL to delete an image")
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
	_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok {
			switch serr.ServiceCode() {
			case azblob.ServiceCodeBlobNotFound:
				return c.Status(404).SendString("Image not found")
			}
		}
		return c.Status(500).SendString("Error deleting blob: " + err.Error())
	}

	return c.Status(200).SendString("Image deleted successfully")
}
