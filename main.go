package main

import (
	"log"
	"os"

	images "github.com/echo-webkom/echo-blob/images"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Use(newAPIKeyAuth)

	app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/api/images", images.HandleGetImageByFilename)
	app.Post("/api/images", images.HandlePostImages)

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	log.Fatal(app.Listen(listenAddr))
}
