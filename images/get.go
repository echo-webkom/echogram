package images

import (
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func HandleGetImageByFilename(c *fiber.Ctx) error {
	filename := c.Query("filename")
	if filename == "" {
		return c.Status(200).SendString("Add ?filename=<filename> to the URL to get an image")
	}

	bm, err := getBlobManager()
	if err != nil {
		return c.Status(500).SendString("Failed to create blob manager")
	}

	data, err := bm.Get(filename)
	if err != nil {
		return c.Status(404).SendString("Image not found")
	}

	c.Type(filepath.Ext(filename))
	return c.Status(200).Send(data)
}
