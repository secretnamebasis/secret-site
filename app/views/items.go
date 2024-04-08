package views

import (
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// HomeData defines the data structure for the home page template
type ItemsData struct {
	Title   string
	Address string
	Items   []models.Item
}

// Home renders the home page
func Items(c *fiber.Ctx) error {
	addr, err := dero.GetWalletAddress()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Retrieve blog posts
	items, err := controllers.AllItems(c)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to retrieve items")
	}
	// Sort by create date, for now
	sort.Slice(
		items,
		func(i, j int) bool {
			return items[i].CreatedAt.String() < items[j].CreatedAt.String()
		},
	)

	// Define data for rendering the template
	data := ItemsData{
		Title:   exports.APP_NAME,
		Address: addr.String(),
		Items:   items,
	}

	// Render the template
	if err := renderTemplate(c, "app/public/items.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
