package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/joho/godotenv"
)

// Config struct to hold configuration parameters
type Server struct {
	Port           int
	Environment    string
	DatabasePath   string
	EnvPath        string
	NodeEndPoint   string
	WalletEndpoint string
}

var (
	Domain = Env(
		"DOMAIN",
	)
	Domainname = "https://" + Domain
)

const ()

var (
	NodeEndPoint      string
	WalletEndpoint    string
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
	// Parse flags
	flag.Parse()

	// Load exports
	Environment = *envFlag
	Port = *portFlag
	DatabaseDir = *dbFlag

	switch Environment {
	case "prod":
		NodeEndPoint = "http://" + Env("DERO_NODE_IP") + ":" + Env("DERO_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env("DERO_WALLET_IP") + ":" + Env("DERO_WALLET_PORT") + "/json_rpc"
		// In production environments, we presuppose DERO mainnet
	case "test":
		Port = 3000
		Domainname = "127.0.0.1"
		NodeEndPoint = "http://" + Env("DERO_SIMULATOR_NODE_IP") + ":" + Env("DERO_SIMULATOR_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env("DERO_SIMULATOR_WALLET_IP") + ":" + Env("DERO_SIMULATOR_WALLET_PORT") + "/json_rpc"
		DatabaseDir = "../app/database/"
		Environment = "test"
		EnvPath = "../.env." + Environment
		dir := "../vendors/derohe/cmd/simulator"
		go func() {
			if err := launchSimulator(dir); err != nil {
				log.Println("Error launching simulator:", err)
			}
		}()
	case "dev":
		Environment = "dev"
		EnvPath = "./.env." + Environment
		Port = 3000
		Domainname = "127.0.0.1"
		NodeEndPoint = "http://" + Env("DERO_SIMULATOR_NODE_IP") + ":" + Env("DERO_SIMULATOR_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env("DERO_SIMULATOR_WALLET_IP") + ":" + Env("DERO_SIMULATOR_WALLET0_PORT") + "/json_rpc"
		// Launch the simulator in the background
		go func() {
			dir := "./vendors/derohe/cmd/simulator"
			if err := launchSimulator(dir); err != nil {
				log.Println("Error launching simulator:", err)
			}
		}()

		time.Sleep(3 * time.Second)
	}

	c := Server{
		Port:         Port,
		Environment:  Environment,
		DatabasePath: DatabaseDir,
		EnvPath:      EnvPath,
		NodeEndPoint: NodeEndPoint,
	}

	return c
}

func launchSimulator(dir string) error {
	cmd := exec.Command("go", "run", ".", "--http-address=0.0.0.0:8081")
	cmd.Dir = dir
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}
	return nil
}
