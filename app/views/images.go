package views

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

func Images(c *fiber.Ctx) error {
	// Extract the image ID from the request URL
	id := c.Params("id")

	// Retrieve the item by ID from the database
	item, err := controllers.GetItemByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Image not found")
	}

	var itemData models.ItemData
	if err := json.Unmarshal(item.Data, &itemData); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	encoded, err := base64.StdEncoding.DecodeString(itemData.Image)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})

	}

	// Detect the content type of the image data
	contentType := http.DetectContentType([]byte(encoded))

	// Set the appropriate content type header
	c.Set(fiber.HeaderContentType, contentType)

	// Send the image data in the response
	return c.Send([]byte(encoded))
}
