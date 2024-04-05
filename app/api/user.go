package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUser creates a new user via HTTP request
func CreateUser(c *fiber.Ctx) error {
	var newUser models.User

	// Parse request body
	if err := c.BodyParser(&newUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate user data
	if err := newUser.Validate(); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Check if user already exists
	if err := checkUserExistence(newUser); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Generate ID and password
	err := controllers.CreateUserRecord(&newUser)

	if err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, "Error creating user")
	}

	return SuccessResponse(c, newUser.Password)
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
