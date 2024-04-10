package views

import (
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
)

func Images(c *fiber.Ctx) error {
	// Extract the image ID from the request URL
	id := c.Params("id")

	// Retrieve the item by ID from the database
	item, err := controllers.GetItemByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Image not found")
	}

	// Decode the base64-encoded image data
	imageData, err := base64.StdEncoding.DecodeString(item.Content.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to decode image data")
	}

	// Set the appropriate content type header
	c.Set(fiber.HeaderContentType, "image/png")

	// Send the image data in the response
	return c.Send(imageData)
}
