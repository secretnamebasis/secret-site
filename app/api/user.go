package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUserHandler creates a new user via HTTP request
func CreateUser(c *fiber.Ctx) error {
	var newUser models.User

	// Parse request body
	if err := c.BodyParser(&newUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate user data
	if err := validateUserData(newUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Check if user already exists
	if err := checkUserExistence(newUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Validate user wallet
	if err := validateWalletAddress(newUser.Wallet); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Generate ID and password
	user, err := createUserRecord(newUser)
	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error creating user")
	}

	return SuccessResponse(c, user.Password)
}

// AllUsers retrieves all users from the database
func AllUsers(c *fiber.Ctx) error {
	users, err := controllers.AllUsers(c)
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
	id := c.Params("id")
	var updatedUser models.User

	// Parse request body
	if err := c.BodyParser(&updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate user data
	if err := validateUserData(updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Check if user already exists
	if err := checkUserExistence(updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Validate user data
	if err := validateWalletAddress(updatedUser.Wallet); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := controllers.UpdateUser(id, updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error updating user")
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

// createUserRecord creates a new user record in the database
func createUserRecord(newUser models.User) (models.User, error) {
	// Generate ID and password
	nextID, _ := controllers.NextUserID()
	password := uuid.New().String()

	user := models.User{
		ID:        nextID,
		User:      newUser.User,
		Wallet:    newUser.Wallet,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := controllers.CreateUserRecord(user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
