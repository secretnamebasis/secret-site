package views

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

func NewAsset(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Get payment ID and message from query parameters
	paymentID := c.Query("payment_id")
	message := c.Query("message")

	// Define data for rendering the template
	data := struct {
		Title     string
		Address   string
		Failed    bool
		PaymentID string
		Message   string
	}{
		Title:     config.Domain,
		Address:   addr.String(),
		Failed:    false,
		PaymentID: paymentID,
		Message:   message,
	}

	// Render the template
	if err := renderTemplate(c, "app/public/asset_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
func SubmitAsset(c *fiber.Ctx) error {
	// Call api.CreateAsset asynchronously
	errCh := make(chan error, 1)
	go func() {
		errCh <- api.CreateAssetOrder(c)
	}()

	// Wait for the response from api.CreateAsset
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
				"address expires",
			),
		):

			// Extract payment ID and message from the response
			var response struct {
				PaymentID string `json:"result"`
				Message   string `json:"message"`
			}
			err := json.Unmarshal(responseBody, &response)
			if err != nil {
				return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
			}

			// Redirect back to /assets/new with payment ID and message as query parameters
			return handleNewAssetSuccess(
				c,
				response.PaymentID,
				response.Message,
			)
		}
	}

	// Redirect to /assets/new upon successful form submission
	return c.Redirect("/assets/new")
}

func handleNewAssetSuccess(c *fiber.Ctx, paymentID, message string) error {
	// Fetch Dero wallet address

	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	data := struct {
		Title     string
		Address   string
		Failed    bool
		PaymentID string
		Message   string
	}{
		Title:     config.Domain,
		Address:   addr.String(),
		Failed:    false,
		PaymentID: paymentID,
		Message:   message,
	}

	// Render the template again with the notice
	if err := renderTemplate(c, "app/public/asset_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
