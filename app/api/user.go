package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUser creates a new user via HTTP request
func CreateUserOrder(c *fiber.Ctx) error {
	order := parseUserData(c)
	if err := controllers.ValidateWalletAddress(order.Wallet); err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}
	if err := controllers.CreateUserRecord(&order); err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}
	return SuccessResponse(
		c,
		"user created",
		&order,
	)
}

// AllUsers retrieves all users from the database
func AllUsers(c *fiber.Ctx) error {
	users, err := controllers.AllUsers()
	if err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			"Error retrieving users",
		)
	}
	return SuccessResponse(
		c,
		"users retrieved",
		users,
	)
}

// UserByID retrieves a user from the database by ID
func UserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := controllers.GetUserByID(id)
	if err != nil {
		return ErrorResponse(
			c,
			fiber.StatusNotFound,
			err.Error(),
		)
	}
	return SuccessResponse(
		c,
		"user retreived",
		user,
	)
}

// UpdateUser updates a user in the database
func UpdateUser(c *fiber.Ctx) error {
	updatedUser := parseUpdatedUserData(c)
	if err := controllers.UpdateUser(updatedUser); err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			err.Error(),
		)
	}
	return SuccessResponse(
		c,
		"user updated",
		nil,
	)
}

// DeleteUser deletes a user from the database
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := controllers.GetUserByID(id); err != nil {
		return ErrorResponse(
			c,
			fiber.StatusNotFound, "User not found",
		)
	}
	if err := controllers.DeleteUser(id); err != nil {
		return ErrorResponse(
			c,
			fiber.StatusInternalServerError,
			"Error deleting user",
		)
	}
	return SuccessResponse(
		c,
		"user deleted",
		nil,
	)
}

//private functions

// parseUserData parses form data or request body to populate the user object
func parseUserData(c *fiber.Ctx) models.JSON_User_Order {
	var order models.JSON_User_Order

	if form, err := c.MultipartForm(); err == nil && form != nil {
		order.Name = form.Value["name"][0]
		order.Wallet = form.Value["wallet"][0]
		order.Password = form.Value["password"][0]
	} else {
		if err := c.BodyParser(&order); err != nil {
			return models.JSON_User_Order{}
		}
		username,
			password,
			_ := getCredentials(c)

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

	return order
}

// parseUpdatedUserData parses request data to update a user
func parseUpdatedUserData(c *fiber.Ctx) models.JSON_User_Order {
	var updatedUser models.JSON_User_Order
	name, password, err := getCredentials(c)
	if err != nil {
		return updatedUser
	}
	if err := c.BodyParser(&updatedUser); err != nil {
		return updatedUser
	}
	if password != "" {
		updatedUser.Password = password
	}
	if name != "" {
		updatedUser.Name = name
	}
	return updatedUser
}
