package main

import (
	"log"
	"os"

	images "github.com/echo-webkom/echo-blob/images"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// TODO: add auth middleware
	// jwt := newAuthMiddleware(os.Getenv("SIGNING_KEY"))

	app.Get("/api/images", images.HandleGetImageByUserId)
	app.Post("/api/images", images.HandlePostImages)
	app.Delete("/api/images", images.HandleDeleteImageByUserId)

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	log.Fatal(app.Listen(listenAddr))
}

func newAuthMiddleware(secret string) fiber.Handler {
	var signKey jwtware.SigningKey
	signKey.Key = []byte(secret)
	return jwtware.New(jwtware.Config{
		SigningKey: signKey,
	})
}
