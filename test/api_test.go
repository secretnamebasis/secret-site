package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
)

// Endpoint configuration
const (
	endpoint             = "http://127.0.0.1:3000/api/users"
	user                 = "secret"
	createAddressFail    = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj"
	updateAddressFail    = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0"
	createAddressSuccess = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj8"
	updateAddressSuccess = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
)

// Configure server settings
var c = config.Server{
	Port: 3000,
	Env:  "testing",
}

// User struct represents the JSON payload for user data
type User struct {
	User   string `json:"user"`
	Wallet string `json:"wallet"`
}

// performAction performs HTTP requests and returns the response status
func performAction(method, url string, data interface{}) (string, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func checkUsers() (string, error) {
	return performAction(
		"GET",
		endpoint,
		nil,
	)
}

func executeTest(t *testing.T, actionFunc func() (string, error), expectedStatus string) {
	actual, err := actionFunc()
	if err != nil {
		t.Fatalf("Error executing action: %v", err)
	}
	var response struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(actual), &response); err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}
	if expectedStatus != response.Status {
		t.Errorf("Expected status: %s, Actual status: %s", expectedStatus, response.Status)
	}
	time.Sleep(1 * time.Second) // Sleep for 1 second
}
func TestActions(t *testing.T) {
	// Start the server and handle shutdown
	a := startServer()

	// Run tests
	runTests(t)

	// Stop the server after tests are done
	stopServer(t, a)

	// Delete the database
	deleteDB()
}

func startServer() *app.App {
	// Delete the database before starting the server
	deleteDB()

	a := app.MakeApp(c)
	go func() {
		if err := a.StartApp(c); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}

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
	}()
	return a
}

func runTests(t *testing.T) {
	// Run tests
	log.Printf("Environment: %s\n", c.Env)

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Get test cases
	testCases := getTestCases()

	// Run your test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fn(t)
		})
	}
}

func stopServer(t *testing.T, a *app.App) {
	// Stop the server after tests are done
	if err := a.StopApp(); err != nil {
		t.Errorf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func deleteDB() {
	dbPath := "./database/"
	err := os.RemoveAll(dbPath)
	if err != nil {
		log.Fatalf("Error deleting database: %s\n", err)
	}
	log.Println("Database deleted successfully")
}

func getTestCases() []struct {
	name string
	fn   func(*testing.T)
} {
	return []struct {
		name string
		fn   func(*testing.T)
	}{
		{
			"CheckUsers",
			func(t *testing.T) { executeTest(t, checkUsers, "success") },
		},
		{
			"Create when invalid",
			createFailTest,
		},
		{
			"Retrieve when user 1 creation fails",
			retrieveUserFailTest,
		},
		// this part is failing
		{
			"Create when valid",
			createSuccessTest,
		},
		{
			"Retrieve when user 1 is successfully created",
			retrieveUserSuccessTest,
		},
		//
		{
			"Update when not valid",
			updateFailTest,
		},
		{
			"Retrieve when user 1 update fails",
			retrieveUserSuccessTest,
		},
		{
			"Update when valid",
			updateSuccessTest,
		},
		{
			"Retrieve when user 1 is successfully updated",
			retrieveUserSuccessTest,
		},
		{
			"Delete when user1 is present",
			deleteUserSuccessTest,
		},
		{
			"Delete when user is not present",
			deleteUserFailTest,
		},
		{
			"Retrieve when user 1 is deleted",
			retrieveUserFailTest,
		},
	}
}

// Functions to perform API CRUD actions

// CREATE
//
// CREATE FAIL
func createFail() (string, error) {
	data := User{
		User:   user,
		Wallet: createAddressFail,
	}
	return performAction("POST", endpoint, data)
}
func createFailTest(t *testing.T) {
	executeTest(t, createFail, "error")
}

// CREATE SUCCESS
func createSuccess() (string, error) {
	data := User{
		User:   user,
		Wallet: createAddressSuccess,
	}
	return performAction("POST", endpoint, data)
}
func createSuccessTest(t *testing.T) {
	executeTest(t, createSuccess, "success")
}

// RETREIVE
func retrieveUser() (string, error) {
	return performAction("GET", fmt.Sprintf("%s/1", endpoint), nil)
}

// RETREIVE SUCCESS
func retrieveUserSuccessTest(t *testing.T) {
	executeTest(t, retrieveUser, "success")
}

// RETREIVE FAIL
func retrieveUserFailTest(t *testing.T) {
	executeTest(t, retrieveUser, "error")
}

// UPDATE
//
// UPDATE FAIL
func updateFail() (string, error) {
	data := User{
		User:   user,
		Wallet: updateAddressFail,
	}
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint), data)
}
func updateFailTest(t *testing.T) {
	executeTest(t, updateFail, "error")
}

// UPDATE SUCCESS
func updateSuccess() (string, error) {
	data := User{
		User:   user,
		Wallet: updateAddressSuccess,
	}
	return performAction("PUT", fmt.Sprintf("%s/1", endpoint), data)
}
func updateSuccessTest(t *testing.T) {
	executeTest(t, updateSuccess, "success")
}

// DELETE
func deleteUser() (string, error) {
	return performAction("DELETE", fmt.Sprintf("%s/1", endpoint), nil)
}

// DELETE SUCCESS
func deleteUserSuccessTest(t *testing.T) {
	executeTest(t, deleteUser, "success")
}

// DELETE FAIL
func deleteUserFailTest(t *testing.T) {
	executeTest(t, deleteUser, "error")
}
