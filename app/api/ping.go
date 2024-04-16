package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"

	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// Ping handles the ping endpoint.
func Ping(c *fiber.Ctx) error {
	d := "pong"

	a, e := dero.GetWalletAddress(config.WalletEndpoint)

	m := "app: " + config.APP_NAME +
		" :: owner: " + a.String()

	s := "success"

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
