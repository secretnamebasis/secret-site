package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse is a common function to generate error responses
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(
		fiber.Map{
			"message": message,
			"status":  "error",
		},
	)
}

// SuccessResponse is a common function to generate success responses
func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(
		fiber.Map{
			"message": message,
			"result":  data,
			"status":  "success",
		},
	)
}
func getCredentials(c *fiber.Ctx) (username, password string, err error) {
	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No Authorization header found
		return "", "", errors.New("no Authorization header found")
	}

	// Extract the username and password from the Authorization header
	decodedCredentials, err := base64.StdEncoding.DecodeString(
		strings.TrimPrefix(
			authHeader,
			"Basic ",
		),
	)

	if err != nil {
		// Error decoding credentials
		return "", "", fmt.Errorf("error decoding credentials: %v", err)
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		// Invalid credentials format
		return "", "", errors.New("invalid credentials format")
	}

	return credentials[0], credentials[1], nil
}
