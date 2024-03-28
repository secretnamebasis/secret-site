package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/exports"
)

// Ping handles the ping endpoint.
func Ping(c *fiber.Ctx) error {
	response := fiber.Map{
		"message": "Welcome to the " + exports.APP_NAME + " API",
		"data":    "pong",
		"status":  "success",
	}
	return c.JSON(response)
}
