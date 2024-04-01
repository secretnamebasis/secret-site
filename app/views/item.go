package views

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

type ItemData struct {
	Title   string
	Address string
	Item    models.Item
}

func Item(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}
	id := c.Params("id")
	item, err := controllers.GetItemByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	data := ItemData{
		Title:   exports.APP_NAME,
		Address: addr.String(),
		Item:    item,
	}

	tmpl, err := template.ParseFiles("app/public/item_detail.html")
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error parsing template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	if err := tmpl.Execute(c.Response().BodyWriter(), data); err != nil {
		// Log the error for debugging
		fmt.Println("Error executing template:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return nil
}
