package main

import (
	"log"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
)

func main() {

	c := config.Initialize()
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
	a.WaitForShutdown()
}
