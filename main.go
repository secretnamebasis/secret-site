package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/exports"
)

// Configure server settings

func main() {
	exports.Env = "prod"
	var c = config.Server{
		Port:         443,
		Env:          exports.Env,
		DatabasePath: "./app/database/" + exports.Env + ".db",
	}
	// Create Fiber app
	a := app.MakeApp(c)
	a.ListenTLS(
		":"+fmt.Sprintf("%d", c.Port),
		"/etc/letsencrypt/live/secretnamebasis.site/cert.pem",
		"/etc/letsencrypt/live/secretnamebasis.site/privkey.pem",
	)
	// Start Fiber app in a separate goroutine
	go func() {
		if err := a.StartApp(c); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for termination signal to stop the server gracefully
	waitForShutdown(a)
}

func waitForShutdown(a *app.App) {
	// Listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server gracefully
	if err := a.StopApp(); err != nil {
		log.Printf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
