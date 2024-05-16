package views

import (
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"

	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// HomeData defines the data structure for the home page template
type UsersData struct {
	Title   string
	Address string
	Users   []models.User
}

func Users(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Retrieve blog posts
	users, err := controllers.AllUsers()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
	// Sort by create date, for now
	sort.Slice(
		users,
		func(i, j int) bool {
			return users[i].Name < users[j].Name
		},
	)

	// Define data for rendering the template
	data := UsersData{
		Title:   config.Domain,
		Address: addr.String(),
		Users:   users,
	}

	// Render the template
	if err := renderTemplate(c, "app/public/users.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
