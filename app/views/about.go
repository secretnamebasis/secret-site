package views

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// Home renders the home page
func About(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	if err := dero.GetWalletAddress(); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
	}{
		Title:   exports.APP_NAME,
		Address: exports.DeroAddress.String(),
	}

	// Render the template
	if err := renderTemplate(c, "app/public/about.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
