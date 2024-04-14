package views

import (
	"encoding/json"
	"net/http"

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
	imageData, err := json.Marshal(item.Data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to decode image data")
	}

	// Detect the content type of the image data
	contentType := http.DetectContentType(imageData)

	// Set the appropriate content type header
	c.Set(fiber.HeaderContentType, contentType)

	// Send the image data in the response
	return c.Send(imageData)
}
