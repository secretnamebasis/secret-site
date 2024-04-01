package api

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

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

	if err := controllers.CreateItemRecord(item); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error creating item")
	}

	return SuccessResponse(c, item)
}

func AllItems(c *fiber.Ctx) error {
	items, err := controllers.AllItems(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving items")
	}
	return SuccessResponse(c, items)
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

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedItem models.Item

	if err := c.BodyParser(&updatedItem); err != nil || updatedItem.Title == "" || updatedItem.Content == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := controllers.UpdateItem(id, updatedItem); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error updating item")
	}

	return SuccessResponse(c, updatedItem)
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
