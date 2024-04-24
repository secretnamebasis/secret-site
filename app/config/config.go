package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
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
	NodeEndpoint   string
	WalletEndpoint string
	AppName        string
	DevAddress     string
	Domain         string
	Domainname     string
}

const ()

var (
	Domain            string
	Domainname        string
	NodeEndpoint      string
	WalletEndpoint    string
	Environment       string
	Port              int
	ProjectDir        = "./"
	EnvPath           string
	DatabaseDir       string
	DeroAddress       *rpc.Address
	DeroAddressResult rpc.GetAddress_Result
	DevAddress        string
	AppName           string
)
var delay = 1 * time.Second

// Config func to get env value from key
func Env(envPath, key string) string {
	if envPath == "" {
		panic(envPath)
	}

	// Load .env file
	err := godotenv.Load(envPath)
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
	case "test":
		Port = 3000
		Environment = "test"
		Domainname = "127.0.0.1"
		EnvPath = "../../.env." + Environment
		NodeEndpoint = "http://" + Env(EnvPath, "DERO_SIMULATOR_NODE_IP") + ":" + Env(EnvPath, "DERO_SIMULATOR_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env(EnvPath, "DERO_SIMULATOR_WALLET_IP") + ":" + Env(EnvPath, "DERO_SIMULATOR_WALLET0_PORT") + "/json_rpc"
		DatabaseDir = "../database/"
		dir := "../../vendors/derohe/cmd/simulator"
		go func() {
			if err := launchSimulator(dir); err != nil {
				log.Println("Error launching simulator:", err)
			}
		}()
		time.Sleep(delay)
	case "dev":
		Environment = "dev"
		EnvPath = "./.env." + Environment
		Port = 3000
		Domainname = "127.0.0.1"
		NodeEndpoint = "http://" + Env(EnvPath, "DERO_SIMULATOR_NODE_IP") + ":" + Env(EnvPath, "DERO_SIMULATOR_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env(EnvPath, "DERO_SIMULATOR_WALLET_IP") + ":" + Env(EnvPath, "DERO_SIMULATOR_WALLET0_PORT") + "/json_rpc"
		// Launch the simulator in the background
		go func() {
			dir := "./vendors/derohe/cmd/simulator"
			if err := launchSimulator(dir); err != nil {
				log.Println("Error launching simulator:", err)
			}
		}()
		time.Sleep(delay)
	case "prod":
		NodeEndpoint = "http://" + Env(EnvPath, "DERO_NODE_IP") + ":" + Env(EnvPath, "DERO_NODE_PORT") + "/json_rpc"
		WalletEndpoint = "http://" + Env(EnvPath, "DERO_WALLET_IP") + ":" + Env(EnvPath, "DERO_WALLET_PORT") + "/json_rpc"
		// In production environments, we presuppose DERO mainnet
	}

	DevAddress = Env(EnvPath, "DEV_ADDRESS")
	AppName = Env(EnvPath, "DOMAIN")

	c := Server{
		Port:         Port,
		Environment:  Environment,
		DatabasePath: DatabaseDir,
		EnvPath:      EnvPath,
		NodeEndpoint: NodeEndpoint,
		DevAddress:   DevAddress,
		AppName:      AppName,
		Domain:       Domain,
		Domainname:   Domainname,
	}

	return c
}
func launchSimulator(dir string) error {
	// Start the simulator in a separate goroutine
	go func() {
		// Start the simulator
		cmd := exec.Command("go", "run", ".", "--http-address=0.0.0.0:8081")
		cmd.Dir = dir
		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start simulator: %v\n", err)
			return
		}

		// Wait for the simulator to finish
		if err := cmd.Wait(); err != nil {
			log.Printf("Error waiting for simulator to finish: %v\n", err)
		}
	}()

	// Check if the simulator is already running and kill it if necessary
	if err := killSimulator(); err != nil {
		return err
	}
	return nil
}

// killSimulator tries to kill the running simulator process
func killSimulator() error {
	// Find the process ID of the simulator
	out, err := exec.Command("pgrep", "simulator").Output()
	if err != nil {
		// If pgrep failed, the simulator is not running
		if strings.Contains(err.Error(), "exit status 1") {
			return nil
		}
		return fmt.Errorf("error finding simulator process ID: %v", err)
	}

	// Parse the process ID
	pidStr := strings.TrimSpace(string(out))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return fmt.Errorf("error parsing simulator process ID: %v", err)
	}

	// Get the process and send a SIGTERM signal to kill it
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("error finding simulator process: %v", err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("error killing simulator process: %v", err)
	}
	log.Printf("Simulator process %d killed", pid)
	return nil
}
