package images

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func HandleGetImageByFilename(c *fiber.Ctx) error {
	filename := c.Query("filename")
	if filename == "" {
		return c.Status(200).SendString("Add ?filename=<filename> to the URL to get an image")
	}

	accountName, accountKey, containerName := getCredentials()
	if accountName == "" || accountKey == "" || containerName == "" {
		return c.Status(500).SendString("Storage account information is not configured")
	}

	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	if err != nil {
		return c.Status(500).SendString("Failed to parse URL")
	}

	creds, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return c.Status(500).SendString("Failed to create credential")
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*URL, pipeline)
	blobURL := containerURL.NewBlobURL(filename)

	ctx := context.Background()
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Println(err)
		return c.Status(404).SendString("Image not found: " + err.Error())
	}

	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 3})
	defer bodyStream.Close()

	data, err := io.ReadAll(bodyStream)
	if err != nil {
		return c.Status(500).SendString("Error reading blob data")
	}

	// _, err = os.Create(filename)
	// if err != nil {
	// 	return c.Status(500).SendString("Error creating file")
	// }

	c.Type(filepath.Ext(filename))
	return c.Status(200).Send(data)
}
