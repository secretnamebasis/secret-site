package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/secretnamebasis/secret-site/app/exports"
)

// Config struct to hold configuration parameters
type Server struct {
	Port         int
	Env          string
	DatabasePath string
	EnvPath      string
}

var (
	Domain = Env(
		"DOMAIN",
	)
	Domainname = "https://" + Domain
)

// Config func to get env value from key
func Env(key string) string {

	// Load .env file
	err := godotenv.Load(exports.EnvPath)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return ""
	}

	return os.Getenv(key)
}

// Configure server settings
var (
	// Define flags for setting the environment and port
	envFlag = flag.String(
		"env",
		"prod", //default
		"environment: dev, prod, test, etc.",
	)

	dbFlag = flag.String(
		"db",
		"./app/database/", //default
		"db location: eg. ./app/database/",
	)

	portFlag = flag.Int(
		"port",
		443, //default
		"server port number",
	)
)

func Initialize() Server {
	// parse flags
	flag.Parse()

	// load exports
	exports.Env = *envFlag
	exports.Port = *portFlag
	exports.DatabaseDir = *dbFlag

	var c = Server{
		Port:         exports.Port,
		Env:          exports.Env,
		DatabasePath: exports.DatabaseDir,
		EnvPath:      exports.EnvPath,
	}
	return c
}
