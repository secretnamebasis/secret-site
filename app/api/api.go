package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/db"
	"github.com/secretnamebasis/secret-site/app/models"
)

// ErrorResponse is a common function to generate error responses
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message, "status": "error"})
}

// SuccessResponse is a common function to generate success responses
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"data": data, "status": "success"})
}

// checkUserExistence checks if a user with the same username or wallet already exists
func checkUserExistence(user models.User) error {
	// Check if user already exists with the same username
	existingUser, err := db.GetUserByUsername(user.User)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser != nil {
		return errors.New("user with the same username already exists")
	}

	// Check if user already exists with the same wallet
	existingUser, err = db.GetUserByWallet(user.Wallet)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser != nil {
		return errors.New("user with the same wallet already exists")
	}

	return nil
}
