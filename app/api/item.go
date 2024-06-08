package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateItemOrder(c *fiber.Ctx) error {
	var order models.JSON_Item_Order
	form, _ := c.MultipartForm()
	// if err != nil {
	// 	return err
	// }
	if form != nil {
		if err := processItemOrderForm(form, &order); err != nil {
			return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	} else {
		// Parse request body into new item
		if err := c.BodyParser(&order); err != nil {
			return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

	}

	// checkout, err := controllers.CreateItemCheckout(&order)
	// if err != nil {
	// 	return ErrorResponse(
	// 		c,
	// 		fiber.StatusInternalServerError,
	// 		err.Error(),
	// 	)
	// }
	// if err := processItemOrderCredentials(
	// 	c,
	// 	&order,
	// ); err != nil {
	// 	return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	// }

	// Create the item record
	item, err := controllers.CreateItemRecord(&order)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Return success response
	return SuccessResponse(c, "item created", &item)
}

func ItemBySCID(c *fiber.Ctx) error {
	scid := c.Params("scid")

	item, err := controllers.GetItemBySCID(scid)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, "item retrieved", item)
}

func AllItems(c *fiber.Ctx) error {
	items, err := controllers.AllItems()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving items")
	}

	return SuccessResponse(c, "item retrieved", items)
}

func UpdateItem(c *fiber.Ctx) error {

	// scid := c.Params("scid")
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

	item, err := controllers.GetItemBySCID(updatedItem.SCID)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return ErrorResponse(c, fiber.StatusNotFound, "Item not found")
		}
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error checking item")
	}

	// update
	if err := controllers.UpdateItem(item.SCID, updatedItem); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, "item updated", &item)
}

func DeleteItem(c *fiber.Ctx) error {
	scid := c.Params("scid")

	// Check if the user exists
	item, err := controllers.GetItemBySCID(scid)
	if err != nil {
		return ErrorResponse(c, fiber.StatusNotFound, "Item not found")
	}
	if err := controllers.DeleteItem(strconv.Itoa(item.ID)); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting item")
	}
	fmt.Printf("%+v", item)
	return SuccessResponse(c, "Item deleted successfully", &item)
}

// private functions
func processItemOrderForm(form *multipart.Form, order *models.JSON_Item_Order) error {
	var imageBase64 string
	if file, ok := form.File["item_data.image"]; ok && len(file) > 0 {
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

	var fileBase64 string

	if file, ok := form.File["item_data.file"]; ok && len(file) > 0 {
		fileFile, err := file[0].Open()
		if err != nil {
			return err
		}
		defer fileFile.Close()

		fileBytes, err := io.ReadAll(fileFile)
		if err != nil {
			return err
		}

		fileBase64 = base64.StdEncoding.EncodeToString(fileBytes)
	}

	// order.User.Name = form.Value["name"][0]
	// order.User.Password = form.Value["password"][0]
	// order.User.Wallet = form.Value["wallet"][0]
	order.Title = form.Value["title"][0]
	order.Description = form.Value["description"][0]
	order.SCID = form.Value["scid"][0]
	order.Image = imageBase64
	order.File = fileBase64 // Assuming order.File is a string field to store the base64 representation of the file
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
