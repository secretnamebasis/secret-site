package views

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

func Files(c *fiber.Ctx) error {
	// Extract the image ID from the request URL
	scid := c.Params("scid")

	// Retrieve the item by ID from the database
	item, err := controllers.GetItemBySCID(scid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}

	var itemData models.ItemData
	if err := json.Unmarshal(item.Data, &itemData); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	// Decode the base64 encoded file data
	decoded, err := base64.StdEncoding.DecodeString(itemData.File)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	// Set the appropriate content type header
	contentType := http.DetectContentType(decoded)
	c.Set(fiber.HeaderContentType, contentType)

	// Set the Content-Disposition header for downloading
	filename := item.FileURL
	c.Set(fiber.HeaderContentDisposition, "attachment; filename="+filename)

	// Send the file data in the response
	return c.Send(decoded)
}
