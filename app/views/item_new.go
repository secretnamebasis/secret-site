package views

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// NewItem renders the new item page
func NewItem(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
		Failed  bool // Add the Failed field
	}{
		Title:   config.Domain,
		Address: addr.String(),
		Failed:  false, // Initially set to false
	}

	// Render the template
	if err := renderTemplate(c, "app/public/item_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// SubmitItem handles the form submission for creating a new item
func SubmitItem(c *fiber.Ctx) error {
	// Call api.CreateItem asynchronously
	errCh := make(chan error, 1)
	go func() {
		errCh <- api.CreateItemOrder(c)
	}()

	// Wait for the response from api.CreateItem
	if err := <-errCh; err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Read the response body
	responseBody := c.Response().Body()

	// Check if the response body is non-empty
	if len(responseBody) > 0 {
		// Switch based on the content of the response body
		switch {
		case bytes.Contains(
			responseBody,
			[]byte(
				"item with the same scid already exists",
			),
		):
			return handleNewItemFailure(
				c,
				"An item with the same SCID already exists. Please choose a different SCID.",
			)
		case bytes.Contains(
			responseBody,
			[]byte(
				"item with the same title already exists",
			),
		):
			return handleNewItemFailure(
				c,
				"An item with the same Title already exists. Please choose a different Title.",
			)
		case bytes.Contains(
			responseBody,
			[]byte(
				"invalid wallet address",
			),
		):
			return handleNewItemFailure(
				c,
				"Invalid wallet address. Please provide a valid DERO wallet address.",
			)

		case bytes.Contains(
			responseBody,
			[]byte(
				"error invalid password",
			),
		):
			return handleNewItemFailure(
				c,
				"Invalid password. Please provide a valid password.",
			)

		case bytes.Contains(
			responseBody,
			[]byte(
				"user does not exist",
			),
		):
			return handleNewItemFailure(
				c,
				"User is not registered. Please register.",
			)
		}
	}

	// Redirect to /items upon successful form submission
	return c.Redirect("/items")
}

// handleNewItemFailure handles the rendering of the registration failure page with a custom message
func handleNewItemFailure(c *fiber.Ctx, message string) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title         string
		Address       string
		Failed        bool   // Flag indicating whether registration failed
		FailedMessage string // Custom failed registration message
	}{
		Title:         config.Domain,
		Address:       addr.String(),
		Failed:        true, // Set to true indicating registration failure
		FailedMessage: message,
	}

	// Render the template again with the notice
	if err := renderTemplate(c, "app/public/item_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
