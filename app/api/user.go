package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUser creates a new user via HTTP request
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	// Parse form data or request body based on content type
	if err := parseUserData(c, &user); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Generate a unique password if the user provided a wallet
	if user.Wallet != "" {
		if err := hasValidWallet(user.Wallet); err != nil {
			return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		user.Password = generateUniquePassword(user.Wallet)
	}

	// Create user record in the database
	if err := controllers.CreateUserRecord(&user); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, user)
}

// parseUserData parses form data or request body to populate the user object
func parseUserData(c *fiber.Ctx, user *models.User) error {
	// Parse form data if available
	if form, err := c.MultipartForm(); err == nil {
		if form != nil {
			user.Name = form.Value["name"][0]
			user.Wallet = form.Value["wallet"][0]
			user.Password = form.Value["password"][0]
			user.Initialize()
		}
	} else {
		// Parse request body
		if err := c.BodyParser(user); err != nil {
			return err
		}

		// Assign default values for missing fields
		username, password, _ := getCredentials(c)
		if user.Name == "" {
			user.Name = username
		}
		if user.Password == "" {
			user.Password = password
		}
		if user.Wallet == "" {
			user.Wallet = c.Params("wallet")
		}
	}
	return nil
}

// generateUniquePassword generates a unique password based on wallet and current timestamp
func generateUniquePassword(wallet string) string {
	// Concatenate wallet and current timestamp
	uniqueData := fmt.Sprintf("%s:%d", wallet, time.Now().UnixNano())
	// Hash the unique data to generate password
	return cryptography.HashString(uniqueData)
}

// AllUsers retrieves all users from the database
func AllUsers(c *fiber.Ctx) error {
	users, err := controllers.AllUsers()
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error retrieving users")
	}
	return SuccessResponse(c, users)
}

// UserByID retrieves a user from the database by ID
func UserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := controllers.GetUserByID(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return SuccessResponse(c, user)
}

// UpdateUser updates a user in the database
func UpdateUser(c *fiber.Ctx) error {
	var updatedUser models.User
	// Get the Authorization header from the request
	name, password, err := getCredentials(c)
	if err != nil {
		// Error getting credentials
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	id := c.Params("id")

	// Parse request body
	if err := c.BodyParser(&updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	intID, err := strconv.Atoi(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	updatedUser.ID = int(intID)

	// creds matter
	if password != "" {
		updatedUser.Password = password
	}

	if name != "" {
		updatedUser.Name = name
	}
	// Validate user data
	if err := updatedUser.Validate(); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := controllers.UpdateUser(id, updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, fiber.Map{"message": "User updated successfully"})
}

// DeleteUser deletes a user from the database
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if the user exists
	_, err := controllers.GetUserByID(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Delete the user
	if err := controllers.DeleteUser(id); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting user")
	}

	return SuccessResponse(c, fiber.Map{"message": "User deleted successfully"})
}
