// site/site.go

package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/db"
	"github.com/secretnamebasis/secret-site/app/routes"
)

func makeApp() *fiber.App {
	app := fiber.New()

	// Initialize the database
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	routes.Draw(app)
	return app
}

func startApp(
	app *fiber.App,
	port int,
) error {
	return app.Listen(
		fmt.Sprintf(
			":%d",
			port,
		),
	)
}
