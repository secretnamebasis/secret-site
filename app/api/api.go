package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
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

// isValidWallet checks if the provided wallet address is valid
func isValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(wallet)
	return err
}

// validateUserData checks if the provided user data is valid
func validateUserData(user models.User) error {
	if user.User == "" || user.Wallet == "" {
		return errors.New("user and wallet fields are required")
	}
	return nil
}

// checkUserExistence checks if a user with the same username or wallet already exists
func checkUserExistence(user models.User) error {
	// Check if user already exists with the same username
	existingUser, err := controllers.GetUserByUsername(user.User)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser != nil {
		return errors.New("user with the same username already exists")
	}

	// Check if user already exists with the same wallet
	existingUser, err = controllers.GetUserByWallet(user.Wallet)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser != nil {
		return errors.New("user with the same wallet already exists")
	}

	return nil
}

// validateWalletAddress checks if the provided wallet address is valid
func validateWalletAddress(wallet string) error {
	if err := isValidWallet(wallet); err != nil {
		return errors.New("invalid wallet address")
	}
	return nil
}
