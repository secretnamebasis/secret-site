package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/api"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

func NewUser(c *fiber.Ctx) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
		Failed  bool // Flag indicating whether registration failed
		Result  string
		Message string
	}{
		Title:   config.Domain,
		Address: addr.String(),
		Failed:  false, // Initially set to false
		Message: "",
		Result:  "",
	}

	// Render the template
	if err := renderTemplate(c, "app/public/user_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// SubmitUser handles the form submission for creating a new user
func SubmitUser(c *fiber.Ctx) error {
	// Call api.CreateUser asynchronously
	errCh := make(chan error, 1)
	go func() {
		errCh <- api.CreateUserOrder(c)
	}()

	// Wait for the response from api.CreateUser
	if err := <-errCh; err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Read the response body
	responseBody := c.Response().Body()

	// Check if the response body is non-empty
	if len(responseBody) > 0 {
		switch {
		case strings.Contains(string(responseBody), "user with the same username already exists"):
			return handleRegistrationFailure(
				c,
				"A user with the same username already exists. Please choose a different username.",
			)
		case strings.Contains(string(responseBody), "user with the same wallet already exists"):
			return handleRegistrationFailure(
				c,
				"A user with the same wallet already exists. Please use a different wallet address.",
			)
		case strings.Contains(string(responseBody), "invalid wallet address"):
			return handleRegistrationFailure(
				c,
				"Invalid wallet address. Please provide a valid DERO wallet address.",
			)
		case strings.Contains(string(responseBody), "user checkout created"):
			// fmt.Printf(string(responseBody))
			var response struct {
				Message string `json:"message"`
				Result  string `json:"result"`
				Status  string `json:"status"`
			}
			err := json.Unmarshal(responseBody, &response)
			if err != nil {
				fmt.Println(err) // Correct way to log the error
				return fiber.NewError(
					http.StatusInternalServerError,
					err.Error(),
				)
			}
			return handleRegistrationSuccess(
				c,
				response.Result,
				response.Message,
			)
		}
	}

	// Redirect to /users upon successful form submission
	return c.Redirect("/users/new")
}

// handleRegistrationFailure handles the rendering of the registration failure page with a custom message
func handleRegistrationFailure(c *fiber.Ctx, message string) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
		Failed  bool // Flag indicating whether registration failed
		Result  string
		Message string // Custom failed registration message
	}{
		Title:   config.Domain,
		Address: addr.String(),
		Failed:  true, // Set to true indicating registration failure
		Result:  "",
		Message: message,
	}

	// Render the template again with the notice
	if err := renderTemplate(c, "app/public/user_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}

// handleRegistrationSuccess handles the rendering of the registration success page with a custom message
func handleRegistrationSuccess(c *fiber.Ctx, result, message string) error {
	// Fetch Dero wallet address
	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to fetch Dero wallet address")
	}

	// Define data for rendering the template
	data := struct {
		Title   string
		Address string
		Failed  bool   // Flag indicating whether registration failed
		Result  string // Custom payment ID
		Message string // Custom success message
	}{
		Title:   config.Domain,
		Address: addr.String(),
		Failed:  false, // Set to false indicating registration success
		Result:  result,
		Message: message,
	}

	// Render the template again with the success notice
	if err := renderTemplate(c, "app/public/user_new.html", data); err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// Set the Content-Type header
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	return nil
}
