package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUser creates a new user in the database
func CreateUser(c *fiber.Ctx) error {
	var newUser models.User

	if err := c.BodyParser(&newUser); err != nil || newUser.User == "" || newUser.Wallet == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	nextID, _ := controllers.NextUserID()
	password := uuid.New().String()

	user := models.User{
		ID:        nextID,
		User:      newUser.User,
		Wallet:    newUser.Wallet,
		Password:  password,
		CreatedAt: time.Now(),
	}

	if err := controllers.CreateUser(user); err != nil {
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

	if err := c.BodyParser(&updatedUser); err != nil || updatedUser.User == "" || updatedUser.Wallet == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := controllers.UpdateUser(id, updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error updating user")
	}

	return SuccessResponse(c, fiber.Map{"message": "User updated successfully"})
}

// DeleteUser deletes a user from the database
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := controllers.DeleteUser(id); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting user")
	}

	return SuccessResponse(c, fiber.Map{"message": "User deleted successfully"})
}
