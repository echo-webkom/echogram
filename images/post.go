package images

import (
	"fmt"

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

	filename := req.Filename

	imageFile, err := req.Open()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to open image")
	}
	defer imageFile.Close()

	bm, err := getBlobManager()
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to create blob manager")
	}

	file := make([]byte, req.Size)
	_, err = imageFile.Read(file)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to read image")
	}

	err = bm.Add(filename, file)
	if err != nil {
		return c.Status(500).SendString("Failed to upload image")
	}

	return c.Status(200).SendString("File uploaded successfully")
}
