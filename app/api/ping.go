package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// Ping handles the ping endpoint.
func Ping(c *fiber.Ctx) error {
	a, e := dero.GetWalletAddress()
	s := "success"
	d := "pong"
	m := "app: " + exports.APP_NAME + " ; " + "owner: " + a.String()

	if e != nil {
		return c.JSON(e)
	}

	r := fiber.Map{
		"message": m,
		"data":    d,
		"status":  s,
	}
	return c.JSON(r)
}
