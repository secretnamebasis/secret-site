package api

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/models"
)

var SECRET = config.Env("SECRET")

func CreateItem(c *fiber.Ctx) error {
	var new models.Item

	if err := c.BodyParser(&new); err != nil || new.Title == "" || new.Content == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	nextID, _ := controllers.NextItemID()
	item := models.Item{
		ID:        nextID,
		Title:     new.Title,
		Content:   new.Content,
		CreatedAt: time.Now(),
	}

	// Encrypt content before storing in the database
	encryptedContent, err := cryptography.EncryptData([]byte(item.Content), SECRET)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error encrypting content")
	}

	// Encode the encrypted content to Base64
	item.Content = base64.StdEncoding.EncodeToString(encryptedContent)

	if err := controllers.CreateItemRecord(item); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error creating item")
	}

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

	// Decode the Base64 encoded content
	decodedBytes, err := base64.StdEncoding.DecodeString(item.Content)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error decoding content")
	}

	// Decrypt the content
	decryptedContent, err := cryptography.DecryptData(decodedBytes, SECRET)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error decrypting content")
	}

	// Set the decrypted content to the item
	item.Content = string(decryptedContent)

	return SuccessResponse(c, item)
}

func AllItems(c *fiber.Ctx) error {
	items, err := controllers.AllItems(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving items")
	}

	// Decrypt the content of each item
	for _, item := range items {
		decodedBytes, err := base64.StdEncoding.DecodeString(item.Content)
		if err != nil {
			return ErrorResponse(c, fiber.StatusInternalServerError, "Error decoding content")
		}

		decryptedContent, err := cryptography.DecryptData(decodedBytes, SECRET)
		if err != nil {
			return ErrorResponse(c, fiber.StatusInternalServerError, "Error decrypting content")
		}

		item.Content = string(decryptedContent)
	}

	return SuccessResponse(c, items)
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedItem models.Item

	if err := c.BodyParser(&updatedItem); err != nil || updatedItem.Title == "" || updatedItem.Content == "" {
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
	encryptedContent, err := cryptography.EncryptData([]byte(updatedItem.Content), SECRET)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error encrypting content")
	}

	// Update the item with the new encrypted content
	item.Content = base64.StdEncoding.EncodeToString(encryptedContent)
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
