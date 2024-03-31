package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/exports"
)

// Configure server settings
var c = config.Server{
	Port: 3000,
	Env:  exports.Env,
}

func main() {
	exports.Env = "prod"
	c.Env = exports.Env
	// Create Fiber app
	a := app.MakeApp(c)

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
