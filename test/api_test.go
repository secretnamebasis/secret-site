package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/models"
)

// SERVER

// pre-conditions
const ( // Endpoint configuration
	endpoint      = "http://127.0.0.1:3000/api"
	routeApiUsers = "/users"
	routeApiItems = "/items"
	user          = "secret"
	dbPath        = "./database/"
)

// test-conditions
var ( // Configure server settings
	c = config.Server{
		Port: 3000,
		Env:  "testing",
	}

	delay = 1 * time.Nanosecond
)

func TestApi(t *testing.T) {
	exports.Env = c.Env
	exports.ProjectDir = "../"
	// Start the server and handle shutdown
	a := startServer()

	// Run tests
	runTests(t)

	// Stop the server after tests are done
	stopServer(t, a)

	// Delete the database
	deleteDB()
}

// test-server
func startServer() *app.App { // start the server
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

func stopServer(t *testing.T, a *app.App) { // stop the server
	// Stop the server after tests are done
	if err := a.StopApp(); err != nil {
		t.Errorf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func deleteDB() {

	err := os.RemoveAll(dbPath)
	if err != nil {
		log.Fatalf("Error deleting database: %s\n", err)
	}
	log.Println("Database deleted successfully")
}

// TEST

// Define testCase type
type testCase struct {
	name string
	fn   func(*testing.T)
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
	time.Sleep(delay) // Sleep for 1 second
}

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func runTests(t *testing.T) {
	// Run tests
	log.Printf("Environment: %s\n", c.Env)

	// Allow some time for the server to start
	time.Sleep(delay)

	// Run user test cases in parallel
	runTestGroup(t, "UserTests", getUserTestCases())

	// Run item test cases in parallel
	runTestGroup(t, "ItemTests", getItemTestCases())
}

func runTestGroup(t *testing.T, groupName string, testCases []testCase) {
	t.Run(
		groupName,
		func(t *testing.T) {
			t.Parallel()
			for _, tc := range testCases {
				tc := tc // Capture range variable
				t.Run(
					tc.name,
					func(t *testing.T) {
						tc.fn(t)
					},
				)
			}
		},
	)
}

// ITEM
var ( // Item Test Data
	//
	// we don't store empty content
	// Fail cases
	failItemCreateData = models.Item{
		Title:   "title",
		Content: "",
	}
	// we don't store empty titles
	failItemUpdateData = models.Item{
		Title:   "",
		Content: "content",
	}
	// Success cases
	successItemCreateData = models.Item{
		Title:   "title",
		Content: "success",
	}
	successItemUpdateData = models.Item{
		Title:   "squirrel",
		Content: "update",
	}
)

func getItemTestCases() []testCase {
	return []testCase{
		{
			"CheckItems",
			func(t *testing.T) { executeTest(t, checkItems, "success") },
		},
		{
			"Create when Item is invalid",
			createItemFailTest,
		},
		{
			"Retrieve when Item 1 creation fails",
			retrieveItemFailTest,
		},
		{
			"Create when Item valid",
			createItemSuccessTest,
		},
		{
			"Retrieve when Item 1 is successfully created",
			retrieveItemSuccessTest,
		},
		{
			"Update when Item is not valid",
			updateItemFailTest,
		},
		{
			"Retrieve when Item 1 update fails",
			retrieveItemSuccessTest,
		},
		{
			"Update when Item is valid",
			updateItemSuccessTest,
		},
		{
			"Retrieve when Item 1 is successfully updated",
			retrieveItemSuccessTest,
		},
		{
			"Delete when Item 1 is present",
			deleteItemSuccessTest,
		},
		{
			"Delete when Item is not present",
			deleteItemFailTest,
		},
		{
			"Retrieve when Item 1 is deleted",
			retrieveItemFailTest,
		},
	}
}

// Functions to perform Item API CRUD actions
func checkItems() (string, error) {
	return performAction(
		"GET",
		endpoint+routeApiItems,
		nil,
	)
}

// CREATE
//
// CREATE FAIL
func createItemFail() (string, error) {

	return performAction("POST", endpoint+routeApiItems, failItemCreateData)
}
func createItemFailTest(t *testing.T) {
	executeTest(t, createItemFail, "error")
}

// CREATE SUCCESS
func createItemSuccess() (string, error) {

	return performAction("POST", endpoint+routeApiItems, successItemCreateData)
}
func createItemSuccessTest(t *testing.T) {
	executeTest(t, createItemSuccess, "success")
}

// RETREIVE
func retrieveItem() (string, error) {
	return performAction("GET", fmt.Sprintf("%s/1", endpoint+routeApiItems), nil)
}

// RETREIVE SUCCESS
func retrieveItemSuccessTest(t *testing.T) {
	executeTest(t, retrieveItem, "success")
}

// RETREIVE FAIL
func retrieveItemFailTest(t *testing.T) {
	executeTest(t, retrieveItem, "error")
}

// UPDATE
//
// UPDATE FAIL
func updateItemFail() (string, error) {

	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiItems), failItemUpdateData)
}
func updateItemFailTest(t *testing.T) {
	executeTest(t, updateItemFail, "error")
}

// UPDATE SUCCESS
func updateItemSuccess() (string, error) {

	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiItems), successItemUpdateData)
}
func updateItemSuccessTest(t *testing.T) {
	executeTest(t, updateItemSuccess, "success")
}

// DELETE
func deleteItem() (string, error) {
	return performAction("DELETE", fmt.Sprintf("%s/1", endpoint+routeApiItems), nil)
}

// DELETE SUCCESS
func deleteItemSuccessTest(t *testing.T) {
	executeTest(t, deleteItem, "success")
}

// DELETE FAIL
func deleteItemFailTest(t *testing.T) {
	executeTest(t, deleteItem, "error")
}

// USER
// test-data
var ( // User Test Data
	// Fail cases
	// resopnse, err := dero.GetEncryptedBalance(address)
	// response.Result.Status != "OK"
	failCreateAddress  = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj"
	failUpdateAddress  = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0"
	failUserCreateData = models.User{
		User:   user,
		Wallet: failCreateAddress,
	}
	failUserUpdateData = models.User{
		User:   user,
		Wallet: failUpdateAddress,
	}
	// Success cases
	successCreateAddress = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj8"
	successUpdateAddress = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"

	successUserCreateData = models.User{
		User:   user,
		Wallet: successCreateAddress,
	}
	successUserUpdateData = models.User{
		User:   user,
		Wallet: successUpdateAddress,
	}
)

func getUserTestCases() []testCase {
	return []testCase{
		{
			"CheckUsers",
			func(t *testing.T) { executeTest(t, checkUsers, "success") },
		},
		{
			"Create when user is invalid",
			createUserFailTest,
		},
		{
			"Retrieve when user 1 creation fails",
			retrieveUserFailTest,
		},
		{
			"Create when user is valid",
			createUserSuccessTest,
		},
		{
			"Retrieve when user 1 is successfully created",
			retrieveUserSuccessTest,
		},
		{
			"Update when user is invalid",
			updateFailTest,
		},
		{
			"Retrieve when user 1 update fails",
			retrieveUserSuccessTest,
		},
		{
			"Update when user is valid",
			updateUserSuccessTest,
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

// Functions to perform User API CRUD actions
func checkUsers() (string, error) {
	return performAction(
		"GET",
		endpoint+routeApiUsers,
		nil,
	)
}

// CREATE
//
// CREATE FAIL
func createUserFail() (string, error) {

	return performAction("POST", endpoint+routeApiUsers, failUserCreateData)
}
func createUserFailTest(t *testing.T) {
	executeTest(t, createUserFail, "error")
}

// CREATE SUCCESS
func createUserSuccess() (string, error) {

	return performAction("POST", endpoint+routeApiUsers, successUserCreateData)
}
func createUserSuccessTest(t *testing.T) {
	executeTest(t, createUserSuccess, "success")
}

// RETREIVE
func retrieveUser() (string, error) {
	return performAction("GET", fmt.Sprintf("%s/1", endpoint+routeApiUsers), nil)
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
func updateUserFail() (string, error) {

	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiUsers), failUserUpdateData)
}
func updateFailTest(t *testing.T) {
	executeTest(t, updateUserFail, "error")
}

// UPDATE SUCCESS
func updateUserSuccess() (string, error) {

	return performAction("PUT", fmt.Sprintf("%s/1", endpoint+routeApiUsers), successUserUpdateData)
}
func updateUserSuccessTest(t *testing.T) {
	executeTest(t, updateUserSuccess, "success")
}

// DELETE
func deleteUser() (string, error) {
	return performAction("DELETE", fmt.Sprintf("%s/1", endpoint+routeApiUsers), nil)
}

// DELETE SUCCESS
func deleteUserSuccessTest(t *testing.T) {
	executeTest(t, deleteUser, "success")
}

// DELETE FAIL
func deleteUserFailTest(t *testing.T) {
	executeTest(t, deleteUser, "error")
}
