// app/site.go

package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
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

	routes.Draw(app)
	return &App{app}
}

// StartApp starts the Fiber application on the specified port
func (a *App) StartApp(c config.Server) error {
	switch exports.Env {
	case "prod":
		return a.ListenTLS(
			":"+fmt.Sprintf("%d", c.Port),
			"/etc/letsencrypt/live/secretnamebasis.site/cert.pem",
			"/etc/letsencrypt/live/secretnamebasis.site/privkey.pem",
		)
	case "dev", "testing":
		return a.Listen(
			fmt.Sprintf(":%d", c.Port),
		)
	default:
		return fmt.Errorf("unsupported environment: %s", exports.Env)
	}
}

// StopApp stops the Fiber application gracefully
func (a *App) StopApp() error {
	return a.Shutdown()
}
