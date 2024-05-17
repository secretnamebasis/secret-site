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
	ServerWallet   string
}

const ()

var (
	Domain         string
	NodeEndpoint   string
	WalletEndpoint string
	Environment    string
	Port           int
	ProjectDir     = "./"
	EnvPath        string
	DatabaseDir    string
	DevAddress     string
	AppName        string
	SimulatorDir   string
)

var (
	UserRegistrationFee  uint64 = 10000
	UserRegistrationPort uint64 = 1337
)

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
	simFlag = flag.Bool(
		"simulator",
		false, //default
		"run in simulator",
	)
)

var delay = 2 * time.Second // Delay for simulator startup

func Initialize() Server {

	// Parse flags
	flag.Parse()

	// Load exports
	Environment = *envFlag
	Port = *portFlag
	DatabaseDir = *dbFlag
	SimulatorDir = "./vendors/derohe/cmd/simulator"

	// Common initialization steps
	switch Environment {
	// prod env
	case "prod":
		initializeForProd()

	// dev env
	case "dev": // assumes mainnet wallet
		initializeForDev()
	case "sim": // assumes simulator wallets
		// the db and .env.sim is different
		// this is to prevent mainnet-testnet conversion errors
		initializeForSim()

	// test env
	case "test":
		initializeForTest()
	}

	if *simFlag {
		launchSimulatorInBackground(SimulatorDir)
	}

	// Create and return the server configuration
	return Server{
		Port:           Port,
		Environment:    Environment,
		DatabasePath:   DatabaseDir,
		EnvPath:        EnvPath,
		NodeEndpoint:   NodeEndpoint,
		WalletEndpoint: WalletEndpoint,
		DevAddress:     DevAddress,
		AppName:        AppName,
		Domain:         Domain,
	}
}

func initializeForTest() {
	Port = 3000
	Environment = "test"
	Domain = "127.0.0.1"
	EnvPath = "../../.env." + Environment
	NodeEndpoint = buildEndpoint("DERO_SIMULATOR_NODE_IP", "DERO_SIMULATOR_NODE_PORT")
	WalletEndpoint = buildEndpoint("DERO_SIMULATOR_WALLET_IP", "DERO_SIMULATOR_WALLET0_PORT")
	DatabaseDir = "../database/"
	SimulatorDir = "../../vendors/derohe/cmd/simulator"
}

func initializeForDev() {
	Environment = "dev"
	EnvPath = "./.env." + Environment
	Port = 3000
	Domain = "127.0.0.1"
	NodeEndpoint = buildEndpoint("DERO_NODE_IP", "DERO_NODE_PORT")
	WalletEndpoint = buildEndpoint("DERO_WALLET_IP", "DERO_WALLET_PORT")
}

func initializeForSim() {
	// this allows for development in non-mainnet conditions, eg simulator
	Environment = "sim"
	EnvPath = "./.env." + Environment
	Port = 3000
	Domain = "127.0.0.1"
	NodeEndpoint = buildEndpoint("DERO_SIMULATOR_NODE_IP", "DERO_SIMULATOR_NODE_PORT")
	WalletEndpoint = buildEndpoint("DERO_SIMULATOR_WALLET_IP", "DERO_SIMULATOR_WALLET0_PORT")

}

func initializeForProd() {
	Environment = "prod"
	EnvPath = "./.env"
	Port = 443
	Domain = Env(EnvPath, "DOMAIN")
	// In production environments, we presuppose DERO mainnet
	NodeEndpoint = buildEndpoint("DERO_NODE_IP", "DERO_NODE_PORT")
	WalletEndpoint = buildEndpoint("DERO_WALLET_IP", "DERO_WALLET_PORT")
}

func buildEndpoint(ipEnvVar, portEnvVar string) string {
	return "http://" + Env(EnvPath, ipEnvVar) + ":" + Env(EnvPath, portEnvVar) + "/json_rpc"
}

func launchSimulatorInBackground(dir string) {
	go func() {
		if err := launchSimulator(dir); err != nil {
			log.Println("Error launching simulator:", err)
		}
	}()
	time.Sleep(delay)
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
