package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/deroproject/derohe/rpc"
	"github.com/joho/godotenv"
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

const ()

var (
	Environment       string
	Port              int
	ProjectDir        = "./"
	EnvPath           = ProjectDir + ".env"
	DatabaseDir       string
	DeroAddress       *rpc.Address
	DeroAddressResult rpc.GetAddress_Result
	DEV_ADDRESS       = Env("DEV_ADDRESS")
	APP_NAME          = Env("DOMAIN")
)

// Config func to get env value from key
func Env(key string) string {

	// Load .env file
	err := godotenv.Load(EnvPath)
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
	Environment = *envFlag
	Port = *portFlag
	DatabaseDir = *dbFlag

	var c = Server{
		Port:         Port,
		Env:          Environment,
		DatabasePath: DatabaseDir,
		EnvPath:      EnvPath,
	}
	return c
}
