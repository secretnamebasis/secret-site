package views

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// NewUser renders the new item page
func NewUser(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
	}{
		Title:   config.APP_NAME,
		Address: addr.String(),
	}

	// Render the template
	if err := renderTemplate(c, "app/public/user_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// // CreateItem handles the form submission for creating a new item
func SubmitUser(c *fiber.Ctx) error {

	api.CreateUser(c)

	// Redirect to /items upon successful form submission
	return c.Redirect("/users")
}
