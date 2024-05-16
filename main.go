package main

import (
	"fmt"
	"log"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/services"
)

func main() {

	c := config.Initialize()
	if c == (config.Server{}) {
		log.Fatalf("Config is empty")
	}

	addr, err := dero.GetWalletAddress(c.WalletEndpoint)
	if err != nil {
		log.Fatalf("Wallet is not loaded")
	}
	c.ServerWallet = addr.String()

	fmt.Printf("WELCOME: %s\n", c.ServerWallet)

	// Create Fiber app
	a := app.MakeApp(c)

	if a == (&app.App{}) {
		log.Fatalf("App is empty")
	}

	// Initialize the database before you run the app
	if err := database.Initialize(c); err != nil {
		log.Fatal(err)
	}
	// Start Fiber app in a separate goroutine
	go func() {
		if err := a.StartApp(c); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()
	if err := services.ProcessCheckouts(c); err != nil {
		log.Fatal(err)
	}

	// Wait for termination signal to stop the server gracefully
	a.WaitForShutdown()
}
