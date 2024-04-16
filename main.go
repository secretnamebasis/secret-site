package main

import (
	"fmt"
	"log"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

func main() {

	c := config.Initialize()

	addr, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		log.Fatalf("Wallet is not loaded")
	}
	fmt.Printf("WELCOME: %s\n", addr.String())
	if c == (config.Server{}) {
		log.Fatalf("Config is empty")
	}

	// Create Fiber app
	a := app.MakeApp(c)

	if a == (&app.App{}) {
		log.Fatalf("App is empty")
	}

	// Start Fiber app in a separate goroutine
	go func() {
		// Initialize the database before you run the app
		if err := database.Initialize(c); err != nil {
			log.Fatal(err)
		}
		if err := a.StartApp(c); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	// Wait for termination signal to stop the server gracefully
	a.WaitForShutdown()
}
