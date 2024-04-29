package views

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// ItemData defines the data structure for the item detail template
type ItemData struct {
	Title       string
	Address     string
	Item        models.Item
	ImageUrl    string
	Image       string
	Description string
}

// Item renders the item detail page
func Item(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Get the item ID from the request parameters
	scid := c.Params("scid")

	// Retrieve the item by ID
	item, err := controllers.GetItemBySCID(scid)
	if err != nil {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": err.Error(),
				"status":  "error",
			},
		)
	}

	var itemData models.ItemData
	if err := json.Unmarshal(item.Data, &itemData); err != nil {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": err.Error(),
				"status":  "error",
			},
		)
	}
	// Define data for rendering the template
	data := ItemData{
		Title:       config.Domain,
		Address:     addr.String(),
		Item:        item,
		ImageUrl:    item.ImageURL,
		Image:       itemData.Image,
		Description: itemData.Description,
	}

	// Render the template using renderTemplate function
	if err := renderTemplate(c, "app/public/item.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
