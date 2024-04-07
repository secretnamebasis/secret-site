package views

import (
	"bytes"
	"html/template"
	"os"

	"github.com/gofiber/fiber/v2"
)

// renderTemplate parses and executes the template with the provided data
func renderTemplate(c *fiber.Ctx, filename string, data interface{}) error {
	// Read the contents of header.html
	headerContent, err := os.ReadFile("app/public/header.html")
	if err != nil {
		return err
	}

	// Parse the header template file
	headerTmpl, err := template.New("header").Parse(string(headerContent))
	if err != nil {
		return err
	}

	// Execute the header template with the provided data
	var header bytes.Buffer
	err = headerTmpl.Execute(&header, data)
	if err != nil {
		return err
	}

	// Include CSS file link in the header
	header.WriteString(`<link rel="stylesheet" href="styles.css">`)

	// Parse the main template file
	mainTmpl, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}

	// Execute the main template file
	var mainContent bytes.Buffer
	err = mainTmpl.Execute(&mainContent, data)
	if err != nil {
		return err
	}

	// Read the contents of footer.html
	footerContent, err := os.ReadFile("app/public/footer.html")
	if err != nil {
		return err
	}

	// Parse the footer template file
	footerTmpl, err := template.New("footer").Parse(string(footerContent))
	if err != nil {
		return err
	}

	// Execute the footer template with the provided data
	var footer bytes.Buffer
	err = footerTmpl.Execute(&footer, data)
	if err != nil {
		return err
	}

	// Write header, main content, and footer contents to the response
	c.Write(header.Bytes())
	c.Write(mainContent.Bytes())
	c.Write(footer.Bytes())

	return nil
}
