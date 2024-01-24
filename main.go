package main

import (
	"log"
	"os"

	images "github.com/echo-webkom/echo-blob/images"
	"github.com/gofiber/fiber/v2"
	cors "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(newAPIKeyAuth)

	app.Get("/api/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "https://echo.uib.no, http://localhost:3000",
			AllowHeaders: "Origin, Content-Type, Accept",
		},
	))

	app.Get("/api/images", images.HandleGetImageByUserId)
	app.Post("/api/images", images.HandlePostImages)
	app.Delete("/api/images", images.HandleDeleteImageByUserId)

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	log.Fatal(app.Listen(listenAddr))
}
