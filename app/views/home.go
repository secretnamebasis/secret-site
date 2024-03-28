package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// HomeData defines the data structure for the home page template
type HomeData struct {
	Title   string
	Address string
	Items   []models.Item
}

// Home renders the home page
func Home(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	if err := dero.Address(); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Retrieve blog posts
	items, err := controllers.AllItems(c)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to retrieve items")
	}

	// Define data for rendering the template
	data := HomeData{
		Title:   exports.APP_NAME,
		Address: exports.DeroAddress.String(),
		Items:   items,
	}

	// Render the template
	if err := renderTemplate(c, "app/public/index.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// renderTemplate parses and executes the template with the provided data
func renderTemplate(c *fiber.Ctx, filename string, data interface{}) error {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	err = tmpl.Execute(c.Response().BodyWriter(), data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	return nil
}
