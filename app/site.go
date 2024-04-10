// app/site.go

package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/routes"
)

// App represents the Fiber application
type App struct {
	*fiber.App
}

// MakeApp creates and initializes a new Fiber application
func MakeApp(c config.Server) *App {
	app := fiber.New(
		fiber.Config{
			AppName:               exports.APP_NAME,
			CaseSensitive:         true,
			DisableStartupMessage: false,
		},
	)

	// Initialize the database
	if err := database.InitDB(c); err != nil {
		log.Fatal(err)
	}

	routes.Draw(app)
	return &App{app}
}

// StartApp starts the Fiber application on the specified port
func (a *App) StartApp(c config.Server) error {
	return a.Listen(
		fmt.Sprintf(":%d", c.Port),
	)
}

// StopApp stops the Fiber application gracefully
func (a *App) StopApp() error {
	return a.Shutdown()
}
