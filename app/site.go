// app/site.go

package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
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
			AppName:               config.Domain,
			CaseSensitive:         true,
			DisableStartupMessage: true,
		},
	)

	routes.Draw(app)
	return &App{app}
}

// StartApp starts the Fiber application on the specified port
func (a *App) StartApp(c config.Server) error {
	switch config.Environment {
	case "prod":
		var cert = "/etc/letsencrypt/live/" + config.Domain + "/cert.pem"
		var privkey = "/etc/letsencrypt/live/" + config.Domain + "/privkey.pem"
		return a.ListenTLS(
			fmt.Sprintf(":%d", c.Port),
			cert,
			privkey,
		)
	case "dev", "sim", "test":
		return a.Listen(
			fmt.Sprintf(":%d", c.Port),
		)
	default:
		return fmt.Errorf("unsupported environment: %s", config.Environment)
	}
}

// StopApp stops the Fiber application gracefully
func (a *App) StopApp() error {
	return a.Shutdown()
}
func (a *App) WaitForShutdown() error {

	// Listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-sigChan

	// Shutdown the server gracefully
	if err := a.StopApp(); err != nil {
		log.Printf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully :)")
	}

	return nil
}
