package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/exports"
)

// Configure server settings

func main() {
	// Define flags for setting the environment and port
	envFlag := flag.String("env", "prod", "environment: dev, prod, test, etc.")
	portFlag := flag.Int("port", 443, "server port number")
	flag.Parse()

	// Set the environment and port from the flag values
	exports.Env = *envFlag
	exports.Port = *portFlag
	var c = config.Server{
		Port:         exports.Port,
		Env:          exports.Env,
		DatabasePath: fmt.Sprintf("./app/database/%s.db", exports.Env),
	}
	// Create Fiber app
	a := app.MakeApp(c)

	// Start Fiber app in a separate goroutine
	go func() {
		// Initialize the database before you run the app
		if err := database.InitDB(c); err != nil {
			log.Fatal(err)
		}
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
