package images

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HandleGetImageByUserId(c *fiber.Ctx) error {
	userId := c.Query("userId")
	if userId == "" {
		return c.Status(400).SendString("Add ?userId=<userId> to the URL to get an image")
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
