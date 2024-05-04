package views

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// ItemData defines the data structure for the item detail template
type UserData struct {
	Title   string
	Address string
	User    models.User
}

// Item renders the item detail page
func User(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Get the item ID from the request parameters
	wallet := c.Params("wallet")

	// Retrieve the item by ID
	user, err := database.GetUserByWallet(wallet)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	// Define data for rendering the template
	data := UserData{
		Title:   config.Domain,
		Address: addr.String(),
		User:    user,
	}

	// Render the template using renderTemplate function
	if err := renderTemplate(c, "app/public/user.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
