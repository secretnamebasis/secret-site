package api

import "github.com/gofiber/fiber/v2"

// ErrorResponse is a common function to generate error responses
func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message, "status": "error"})
}

// SuccessResponse is a common function to generate success responses
func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{"data": data, "status": "success"})
}
