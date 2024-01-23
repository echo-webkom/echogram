package images

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofiber/fiber/v2"
)

func HandlePostImages(c *fiber.Ctx) error {
	req, err := c.FormFile("image")
	if err != nil {
		fmt.Println("ERROR", err)
		return c.Status(500).SendString("Failed to decode image")
	}

	if req.Size == 0 {
		return c.Status(400).SendString("File is empty")
	}

	if req.Size > 1024*1024*4 {
		return c.Status(400).SendString("File is too big. Limit is 4MB")
	}

	if !validImageType(req.Filename) {
		return c.Status(400).SendString("Invalid image type. Valid types are .jpg, .jpeg, .png, and .gif")
	}

	imageFile, err := req.Open()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to open image")
	}
	defer imageFile.Close()

	userId := c.Get("User-ID")
	if userId == "" {
		return c.Status(400).SendString("User ID is missing")
	}

	ext := filepath.Ext(req.Filename)
	filename := userId + ext

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to create file")
	}
	defer file.Close()

	_, err = io.Copy(file, imageFile)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to copy image")
	}

	accountName, accountKey, containerName := getCredentials()
	if accountName == "" || accountKey == "" || containerName == "" {
		return c.Status(500).SendString("Storage account information is not configured")
	}

	creds, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to create credentials")
	}

	pipeline := azblob.NewPipeline(creds, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, pipeline)

	ctx := context.Background()
	blobURL := containerURL.NewBlockBlobURL(filename)

	// if user already has an image, delete it and upload the new one
	listBlobs, err := containerURL.ListBlobsFlatSegment(ctx, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to list blobs: " + err.Error())
	}

	for _, blob := range listBlobs.Segment.BlobItems {
		if strings.HasPrefix(blob.Name, userId+".") && validImageType(blob.Name) {
			_, err = containerURL.NewBlobURL(blob.Name).Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
			if err != nil {
				fmt.Println(err)
				return c.Status(500).SendString("Failed to delete old image: " + err.Error())
			}
		}
	}

	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		fmt.Println("Failed to upload file: ", err)
		return c.Status(500).SendString("Failed to upload file")
	}

	return c.Status(200).SendString("File uploaded successfully")
}
