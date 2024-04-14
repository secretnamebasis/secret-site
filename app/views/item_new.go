package views

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// NewItem renders the new item page
func NewItem(c *fiber.Ctx) error {
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
	if err := renderTemplate(c, "app/public/item_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// CreateItem handles the form submission for creating a new item
func SubmitItem(c *fiber.Ctx) error {

	// Parse form data
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// Extract form values

	// Check if an image file was uploaded
	var imageBase64 string
	imageBase64 = ""
	if file, ok := form.File["itemdata.image"]; ok && len(file) > 0 {
		// Check MIME type of the uploaded file
		imageFile, err := file[0].Open()
		if err != nil {
			return err
		}
		defer imageFile.Close()

		// Read the first 512 bytes to detect the MIME type
		buffer := make([]byte, 512)
		_, err = imageFile.Read(buffer)
		if err != nil {
			return err
		}

		// Reset the file offset to start
		_, err = imageFile.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}

		// Detect MIME type
		mimeType := http.DetectContentType(buffer)
		if !strings.HasPrefix(mimeType, "image/") {
			// Return error if the uploaded file is not an image
			return errors.New("invalid file format, please upload an image")
		}

		// Read file contents
		imageBytes, err := io.ReadAll(imageFile)
		if err != nil {
			return err
		}

		// Encode image bytes as base64
		imageBase64 = base64.StdEncoding.EncodeToString(imageBytes)
	}
	var item models.JSON_Item_Order

	item.Title = form.Value["title"][0]
	// Convert ItemData into bytes
	item.Description = form.Value["description"][0]
	item.Image = imageBase64

	api.CreateItem(c)

	// Redirect to /items upon successful form submission
	return c.Redirect("/items")
}
