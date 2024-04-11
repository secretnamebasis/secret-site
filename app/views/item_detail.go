package views

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"

	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// ItemData defines the data structure for the item detail template
type ItemData struct {
	Title    string
	Address  string
	Item     models.Item
	ImageUrl string
}

// Item renders the item detail page
func Item(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Get the item ID from the request parameters
	id := c.Params("id")

	// Retrieve the item by ID
	item, err := controllers.GetItemByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	// Define data for rendering the template
	data := ItemData{
		Title:    config.APP_NAME,
		Address:  addr.String(),
		Item:     item,
		ImageUrl: item.Content.ImageURL,
	}

	// Render the template using renderTemplate function
	if err := renderTemplate(c, "app/public/item_detail.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
