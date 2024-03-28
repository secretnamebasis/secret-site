// views/notfound.go
package views

import "github.com/gofiber/fiber/v2"

func NotFound(c *fiber.Ctx) error {
	return c.Status(
		fiber.StatusNotFound,
	).SendString(
		"404 Not Found",
	)
}
