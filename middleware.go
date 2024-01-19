package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

func newAPIKeyAuth(c *fiber.Ctx) error {
	if os.Getenv("ENV") == "dev" {
		return c.Next()
	}

	return keyauth.New(keyauth.Config{
		KeyLookup: "header:Authorization",
		Validator: validateAPIKey,
	})(c)
}

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	if key == "1234" {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
