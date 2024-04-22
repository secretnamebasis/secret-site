package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUser creates a new user via HTTP request
func CreateUser(c *fiber.Ctx) error {
	var order models.JSON_User_Order

	// Parse form data or request body based on content type
	if err := parseUserData(c, &order); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if order.Wallet == "" {
		return ErrorResponse(c, fiber.StatusBadRequest, "No Wallet Found")
	}

	// Create user record in the database
	if err := controllers.CreateUserRecord(&order); err != nil {
		return ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return SuccessResponse(c, order)
}

// parseUserData parses form data or request body to populate the user object
func parseUserData(c *fiber.Ctx, order *models.JSON_User_Order) error {
	// Parse form data if available
	if form, err := c.MultipartForm(); err == nil {
		if form != nil {
			order.Name = form.Value["name"][0]
			order.Wallet = form.Value["wallet"][0]
			order.Password = form.Value["password"][0]
		}
	} else {
		// Parse request body
		if err := c.BodyParser(order); err != nil {
			return err
		}

		// Assign default values for missing fields
		username, password, _ := getCredentials(c)
		if order.Name == "" {
			order.Name = username
		}
		if order.Password == "" {
			order.Password = password
		}
		if order.Wallet == "" {
			order.Wallet = c.Params("wallet")
		}
	}
	return nil
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
	intID, err := strconv.Atoi(id)
	if err != nil {
		return ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}
	// Parse request body
	if err := c.BodyParser(&updatedUser); err != nil {
		return ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updatedUser.ID = intID

	// creds matter
	if password != "" {
		updatedUser.Password = []byte(password)
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
