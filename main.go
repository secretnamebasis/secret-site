package main

import (
	"log"

	"github.com/secretnamebasis/secret-site/app/config"
)

func main() {

	app := makeApp()

	config := config.Server{
		Port: 3000,
	}

	if err := startApp(
		app,
		config.Port,
	); err != nil {

		log.Printf(
			"Error starting server: %s\n",
			err,
		)

	}
}
