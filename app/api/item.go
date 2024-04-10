package api

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateItem(c *fiber.Ctx) error {
	// Parse request body into new item
	var new models.Item
	if err := c.BodyParser(&new); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Validate the new item
	if err := new.Validate(); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid item: "+err.Error())
	}

	// Get the next item ID
	nextID, _ := controllers.NextItemID()

	// Create a new item with the parsed data
	item := models.InitializeItem(
		nextID,
		new.Title,
		new.Content.Description,
		new.Content.Image,
	)

	// Create the item record
	if err := controllers.CreateItemRecord(item); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error creating item")
	}

	// Return success response
	return SuccessResponse(c, item)
}

func ItemByID(c *fiber.Ctx) error {
	id := c.Params("id")

	item, err := controllers.GetItemByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return SuccessResponse(c, item)
}

func AllItems(c *fiber.Ctx) error {
	items, err := controllers.AllItems()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving items")
	}

	return SuccessResponse(c, items)
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedItem models.Item

	if err := c.BodyParser(&updatedItem); err != nil || updatedItem.Title == "" || updatedItem.Content.Description == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Check if the item exists
	item, err := controllers.GetItemByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ErrorResponse(c, fiber.StatusNotFound, "Item not found")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error checking item")
	}

	// Encrypt the new content
	encryptedContent, err := cryptography.EncryptData([]byte(updatedItem.Content.Description), config.Env("SECRET"))
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error encrypting content")
	}

	// Update the item with the new encrypted content
	item.Content.Description = base64.StdEncoding.EncodeToString(encryptedContent)
	if err := controllers.UpdateItem(id, item); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error updating item")
	}

	return SuccessResponse(c, item)
}

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	// Check if the user exists
	_, err := controllers.GetItemByID(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}
	err = controllers.DeleteItem(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting item")
	}

	return SuccessResponse(c, "Item deleted successfully")
}
