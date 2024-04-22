package api

import (
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

func processItemOrderForm(form *multipart.Form, order *models.JSON_Item_Order) error {
	var imageBase64 string
	imageBase64 = ""
	if file, ok := form.File["itemdata.image"]; ok && len(file) > 0 {
		imageFile, err := file[0].Open()
		if err != nil {
			return err
		}
		defer imageFile.Close()

		buffer := make([]byte, 512)
		_, err = imageFile.Read(buffer)
		if err != nil {
			return err
		}

		_, err = imageFile.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		mimeType := http.DetectContentType(buffer)
		if !strings.HasPrefix(mimeType, "image/") {
			return errors.New("invalid file format, please upload an image")
		}

		imageBytes, err := io.ReadAll(imageFile)
		if err != nil {
			return err
		}

		imageBase64 = base64.StdEncoding.EncodeToString(imageBytes)
	}

	order.Title = form.Value["title"][0]
	order.Description = form.Value["description"][0]
	order.Image = imageBase64
	order.User.Wallet = form.Value["wallet"][0]
	order.User.Name = form.Value["name"][0]
	order.User.Password = form.Value["password"][0]
	return nil
}
func processItemOrderCredentials(c *fiber.Ctx, order *models.JSON_Item_Order) error {

	name, pass, err := getCredentials(c)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if order.User.Name == "" {
		order.User.Name = name
	}
	if order.User.Password == "" {
		order.User.Password = pass
	}
	if order.User.Wallet == "" {
		user, err := controllers.GetUserByName(order.User.Name)
		if err != nil {
			return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		}
		order.User.Wallet = user.Wallet
	}
	return nil
}
func CreateItem(c *fiber.Ctx) error {
	var order models.JSON_Item_Order
	form, _ := c.MultipartForm()
	// if err != nil {
	// 	return err
	// }
	if form != nil {
		processItemOrderForm(form, &order)
	}

	// Parse request body into new item
	if err := c.BodyParser(&order); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := processItemOrderCredentials(c, &order); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Create the item record
	item, err := controllers.CreateItemRecord(&order)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
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
	var updatedItem models.JSON_Item_Order

	if err := c.BodyParser(&updatedItem); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := processItemOrderCredentials(c, &updatedItem); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	if err := updatedItem.Validate(); err != nil {
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

	if err := controllers.UpdateItem(id, updatedItem); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, &item)
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
