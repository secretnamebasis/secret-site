package views

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// HomeData defines the data structure for the home page template
type HomeData struct {
	Title   string
	Address string
}

// Home renders the home page
func Home(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := HomeData{
		Title:   config.APP_NAME,
		Address: addr.String(),
	}

	// Render the template
	if err := renderTemplate(c, "app/public/index.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
