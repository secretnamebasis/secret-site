package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/models"
)

// # API_TEST

// TABLE OF CONTENTS:
// - INTRO
// - SERVER
// - CONFIG
// - RESPONSE
// - DELAY
// - MAIN
// - EXECUTION
// - ACTION
// - CASES
// - TESTS
//	- ITEMS
//	 - CREATE
//	 - RETRIEVE
//	 - UPDATE
//	 - DELETE
//	- USERS
//	 - CREATE
//	 - RETRIEVE
//	 - UPDATE
//	 - DELETE
// - DATA

// INTRO:
// Testing is often considered the bulwark and bane
// of web devs around. The secret-site takes a
// "does the API work?" approach to development.

// The core reasonings to "does the API work" method:
// - delimits the number of tests that we need to maintain
// - simplifies the development scope to one focus: the API

// With this principle in mind, the contemporary web site is
// in the business of serving its code with these pre-requisites:
// - develop the API as an owner
// - use the API as a customer

// By understanding our relationship to the API, as an owner
// and as a customer, the project takes on both simplicity
// as well as scalability - we know what we need our API to do.

// SERVER
const // Endpoint configuration
(
	endpoint      = "http://127.0.0.1:3000/api"
	routeApiUsers = "/users/"
	routeApiItems = "/items/"
	user          = "secret"
	pass          = "pass"
	user2         = pass
	ID            = "1"
)

func // CONFIG
configServer() config.Server {
	config.Environment = "test"
	config.DatabaseDir = "../app/database/"
	config.EnvPath = "./.env"
	c := config.Server{
		Port:         3000, // production runs with :443 for TLS connections
		Env:          config.Environment,
		DatabasePath: config.DatabaseDir,
		EnvPath:      config.EnvPath,
	}
	return c
}

type // RESPONSE
response struct {
	Result interface{} `json:"result"`
	Status string      `json:"status"`
}

var // DELAY
delay = 0 * time.Nanosecond

func briefPause() { time.Sleep(1 * time.Millisecond) }

// MAIN
func TestAPI(t *testing.T) {
	// Check if testing framework is empty
	checkTestingFramework(t)

	// Config server
	cfg := configServer()

	// Check if config is empty
	checkConfig(cfg)

	// Initialize the database with configs
	initializeDatabase(cfg)

	// Start the server as an app
	app := startServer(t, cfg)

	// Check if app is empty
	checkApp(app)

	// Let the server turn on
	briefPause()

	//
	// if err := launchSimulator(); err != nil {
	// 	log.Fatal(err)
	// }

	// // Run tests with configs
	runTests(t, cfg)

	// Stop the server after tests are done
	stopServer(t, app)

	// Delete the database
	deleteTestDB(cfg)
}

func checkTestingFramework(t *testing.T) {
	if t == nil || t == (&testing.T{}) {
		log.Fatal("Testing framework is empty")
	}
}

func checkConfig(cfg config.Server) {
	if cfg == (config.Server{}) {
		log.Fatal("Configuration is empty")
	}
}

func initializeDatabase(cfg config.Server) {
	if err := database.Initialize(cfg); err != nil {
		log.Fatal(err)
	}
}

func checkApp(appInstance *app.App) {
	if appInstance == nil || appInstance == (&app.App{}) {
		log.Fatal("App is empty")
	}
}

func startServer(t *testing.T, c config.Server) *app.App { // start the server

	a := app.MakeApp(c)
	if a == (&app.App{}) {
		log.Fatalf("App is empty")
	}
	go func() {
		if err := a.StartApp(c); err != nil {
			t.Errorf("Error starting server: %s\n", err)
		}
	}()
	return a
}

func launchSimulator() error {
	dir := "../vendors/derohe/cmd/simulator"
	// Command to execute the simulator.go file
	cmd := exec.Command("go", "run", ".", "--http-address=0.0.0.0:8081")

	// Set the working directory
	cmd.Dir = dir

	// Create a pipe to capture the command output
	// stdoutPipe, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return fmt.Errorf("error creating stdout pipe: %v", err)
	// }

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Create a scanner to read from the pipe
	// scanner := bufio.NewScanner(stdoutPipe)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	fmt.Println(line) // Do something with the output line
	// }

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}

	return nil
}
func runTests(t *testing.T, c config.Server) { // run tests
	log.Printf("Environment: %s\n", c.Env)
	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, tc.fn)
	}
}
func stopServer(t *testing.T, a *app.App) { // stop the server
	// Stop the server after tests are done
	if err := a.StopApp(); err != nil {
		t.Errorf("Error stopping server: %s\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
func deleteTestDB(c config.Server) { // clean up
	err := os.Remove(c.DatabasePath + c.Env + ".db")
	if err != nil {
		log.Fatalf("Error deleting database: %s\n", err)
	}
	log.Println("Database deleted successfully")
}

func // EXECUTION
execute(t *testing.T, actionFunc func() (string, error), validateFunc func(string) bool) {
	// Execute the action function
	responseBody, err := actionFunc()
	if err != nil {
		t.Fatalf("Error executing action: %v", err)
	}

	// Perform custom validation
	if !validateFunc(responseBody) {
		t.Errorf("Validation failed for response: %s", responseBody)
	}
	// log.Printf("%s", resp)
	// Sleep for 0 nanosecond
	time.Sleep(delay)
}
func // ACTION
action(method, url string, data interface{}) (string, error) {
	// Marshal data into JSON payload
	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON data: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add basic authentication
	req.SetBasicAuth(user, pass)

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

// CASES
type testCase struct { // structure test cases
	name string
	fn   func(*testing.T)
}

var (
	// array of test cases

	testCases = []testCase{

		{
			// but more valuable to me is not the items
			// the most value item that I want is the db
			// we are trying to determine if we are a user
			// and we expect that we aren't a user among users

			// what's suprising is that because the application
			// has evolved, so has this point in the test of the
			// api. I think it is suprising that this point would
			// change the way that it has.
			"Retrieve error when Users do not exist",
			retreiveUsersTestNoDataSuccess,
		},
		// CREATE TESTS
		{
			// we expect that when there is nothing present in the database,
			// there will be an error when retreiving all items
			"Retrieve error when Items do not exist",
			retrieveItemsTestNoDataSuccess,
		},
		{
			"Create error when Item 1 is invalid",
			createItemTestFail,
		},
		{
			"Retrieve error when Item 1 does not exist",
			retrieveItemTestFail,
		},
		{
			"Create error when User 1 is invalid",
			createUserTestFail,
		},
		{
			"Retrieve error when User 1 does not exist",
			retrieveUserTestFail,
		},
		{
			"Create success when User 1 is valid",
			createUserTestSuccess,
		},
		{
			"Create fail when User 1 already exists",
			createUserTestDuplicateDataFail,
		},
		{
			"Create success when User 2 is valid",
			createUserTestSecondSuccess,
		},
		{
			"Retrieve success when User 1 exists",
			retrieveUserTestSuccess,
		},
		{
			"Create success when Item 1 is valid",
			createItemTestSuccess,
		},
		{
			// we expect that the user is already created
			// otherwise, anyone can make items before there is a user
			"Create error when Item 1 already exists",
			createItemTestDuplicateFail,
		},
		{
			"Retrieve success Item 1 when Item 1 exists",
			retrieveItemTestSuccess,
		},
		{
			"Retrieve success when Items exist",
			retrieveItemsTestSuccess,
		},
		{
			// the user runs into problems.
			"Update error when Item 1 is invalid",
			updateItemTestFail,
		},
		{
			"Retrieve success when Item 1 update fails",
			retrieveItemTestSuccess,
		},
		{
			"Update success when Item 1 is valid",
			updateItemTestSuccess,
		},
		{
			"Retrieve success when Item 1 is updated",
			retrieveItemTestUpdateSuccess,
		},
		{
			"Delete success when Item 1 exisits",
			deleteItemTestSuccess,
		},
		{
			"Delete error when Item 1 does not exist",
			deleteItemTestFail,
		},
		{
			"Retrieve error when Item 1 is deleted",
			retrieveItemTestFail,
		},

		{
			"Update error when User 1 is invalid",
			updateUserTestFail,
		},
		{
			"Retrieve when user 1 update fails",
			retrieveUserTestSuccess,
		},
		{
			"Retrieve sucess when Items exist",
			retreiveUsersTestSuccess,
		},
		{
			"Update success when User 1 is valid",
			updateUserTestSuccess,
		},
		{
			"Retrieve success when User 1 is updated",
			retrieveUserTestSuccess,
		},
		{
			"Delete success when User 1 exisits",
			deleteUserTestSuccess,
		},
		{
			"Delete error when User 1 does not exist",
			deleteUserTestFail,
		},
		{
			"Retrieve error when User 1 is deleted",
			retrieveUserTestFail,
		},
	}
)

// TESTS
func // Retrieve All Items
checkItems() (string, error) {
	return action(
		"GET",
		endpoint+routeApiItems,
		nil,
	)
}
func retrieveItemsTestFail(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the status is "error"
		return resp.Status == "error"
	}

	// Execute the test with custom validation
	execute(t, checkItems, validateFunc)
}
func retrieveItemsTestNoDataSuccess(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the data is not nil and the status is "success"
		return resp.Result == nil && resp.Status == "success"
	}

	// Execute the test with custom validation
	execute(t, checkItems, validateFunc)
}
func retrieveItemsTestSuccess(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the data is not nil and the status is "success"
		return resp.Result != nil && resp.Status == "success"
	}

	// Execute the test with custom validation
	execute(t, checkItems, validateFunc)
}

func // CREATE
createItem(createData interface{}) func() (string, error) {
	return func() (string, error) {
		return action(
			"POST",
			endpoint+routeApiItems,
			createData,
		)
	}
}
func // CREATE FAIL
createItemTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}
		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, createItem(failItemCreateData), validateFunc)
}

// CREATE SUCCESS
func createItemTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Check if the status is "success"
		if resp.Status != "success" {
			t.Errorf("Expected status to be 'success', got '%s'", resp.Status)
			return false
		}

		// Extract the result field
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Errorf("Unexpected type for response data")
			return false
		}

		// Perform validations based on the result data
		// Validate ID
		id, idOK := result["id"].(float64)
		if !idOK || int(id) != 1 {
			t.Errorf("Expected ID to be 1, got %v", id)
			return false
		}

		// Validate title
		title, titleOK := result["title"].(string)
		if !titleOK || title != "First Post" {
			t.Errorf("Expected title to be 'First Post', got '%s'", title)
			return false
		}

		// Validate data
		encodedData, dataOK := result["data"].(string)
		if !dataOK {
			t.Errorf("Expected 'data' field to be present")
			return false
		}
		fmt.Printf("ENCODED DATA: %s\n", encodedData)
		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		fmt.Printf("DECODED DATA: %s\n", decodedData)
		// Decrypt the data
		decryptedData, err := cryptography.DecryptData(decodedData, config.Env("SECRET"))
		if err != nil {
			t.Errorf("Error decrypting data: %v", err)
			return false
		}
		fmt.Printf("DECRYPTED DATA: %s\n", decryptedData)
		// Unmarshal decrypted data into a struct
		var itemData models.ItemData
		if err := json.Unmarshal(decryptedData, &itemData); err != nil {
			t.Errorf("Error unmarshaling decrypted data: %v", err)
			return false
		}

		// Validate description
		if itemData.Description != "love you Joyce" {
			t.Errorf("Expected description to be 'love you Joyce', got '%s'", itemData.Description)
			return false
		}

		return true
	}

	execute(t, createItem(successItemCreateData), validateFunc)
}

// CREATE SUCCESS
func createItemTestDuplicateFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}
		return resp.Status == "error"
	}

	execute(t, createItem(successItemCreateData), validateFunc)
}

func // RETREIVE
retrieveItem() (string, error) {
	return action(
		"GET",
		endpoint+routeApiItems+ID,
		nil,
	)
}
func // RETRIEVE SUCCESS
retrieveItemTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Check if the status is "success"
		if resp.Status != "success" {
			t.Errorf("Expected status to be 'success', got '%s'", resp.Status)
			return false
		}

		// Extract the result field
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Errorf("Unexpected type for response data")
			return false
		}

		// Perform validations based on the result data
		// Validate ID
		id, idOK := result["id"].(float64)
		if !idOK || int(id) != 1 {
			t.Errorf("Expected ID to be 1, got %v", id)
			return false
		}

		// Validate title
		title, titleOK := result["title"].(string)
		if !titleOK || title != "First Post" {
			t.Errorf("Expected title to be 'First Post', got '%s'", title)
			return false
		}

		// Validate data
		encodedData, dataOK := result["data"].(string)
		if !dataOK {
			t.Errorf("Expected 'data' field to be present")
			return false
		}
		fmt.Printf("ENCODED DATA: %s\n", encodedData)

		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		fmt.Printf("DECODED DATA: %s\n", decodedData)
		// Unmarshal decrypted data into a struct
		var itemData models.ItemData
		if err := json.Unmarshal(decodedData, &itemData); err != nil {
			t.Errorf("Error unmarshaling decrypted data: %v", err)
			return false
		}

		// Validate description
		if itemData.Description != "love you Joyce" {
			t.Errorf("Expected description to be 'love you Joyce', got '%s'", itemData.Description)
			return false
		}

		return true
	}

	execute(t, retrieveItem, validateFunc)
}
func // RETRIEVE SUCCESS
retrieveItemTestUpdateSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Check if the status is "success"
		if resp.Status != "success" {
			t.Errorf("Expected status to be 'success', got '%s'", resp.Status)
			return false
		}

		// Extract the result field
		result, ok := resp.Result.(map[string]interface{})
		if !ok {
			t.Errorf("Unexpected type for response data")
			return false
		}

		// Perform validations based on the result data
		// Validate ID
		id, idOK := result["id"].(float64)
		if !idOK || int(id) != 1 {
			t.Errorf("Expected ID to be 1, got %v", id)
			return false
		}

		// Validate title
		title, titleOK := result["title"].(string)
		if !titleOK || title != "squirrel" {
			t.Errorf("Expected title to be 'squirrel', got '%s'", title)
			return false
		}

		// Validate data
		encodedData, dataOK := result["data"].(string)
		if !dataOK {
			t.Errorf("Expected 'data' field to be present")
			return false
		}
		fmt.Printf("ENCODED DATA: %s\n", encodedData)

		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		fmt.Printf("DECODED DATA: %s\n", decodedData)
		// Unmarshal decrypted data into a struct
		var itemData models.ItemData
		if err := json.Unmarshal(decodedData, &itemData); err != nil {
			t.Errorf("Error unmarshaling decrypted data: %v", err)
			return false
		}

		// Validate description
		if itemData.Description != "Some words to drive you nuts" {
			t.Errorf("Expected description to be 'Some words to drive you nuts', got '%s'", itemData.Description)
			return false
		}

		return true
	}

	execute(t, retrieveItem, validateFunc)
}
func // RETRIEVE FAIL
retrieveItemTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, retrieveItem, validateFunc)
}
func // UPDATE
updateItem(updateData interface{}) func() (string, error) {
	return func() (string, error) {
		return action(
			"PUT",
			endpoint+routeApiItems+ID,
			updateData,
		)
	}
}
func // UPDATE FAIL
updateItemTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, updateItem(failItemUpdateData), validateFunc)
}
func // UPDATE SUCCESS
updateItemTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, updateItem(successItemUpdateData), validateFunc)
}
func // DELETE
deleteItem() (string, error) {
	return action(
		"DELETE",
		endpoint+routeApiItems+ID,
		nil,
	)
}
func // DELETE FAIL
deleteItemTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, deleteItem, validateFunc)
}
func // DELETE SUCCESS
deleteItemTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, deleteItem, validateFunc)
}
func // USER
retreiveUsers() (string, error) {
	return action(
		"GET",
		endpoint+routeApiUsers,
		nil,
	)
}
func retreiveUsersTestFail(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the status is "error"
		return resp.Result == nil && resp.Status == "error"
	}

	// Execute the test with custom validation
	execute(t, retreiveUsers, validateFunc)
}
func retreiveUsersTestNoDataSuccess(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the data is not nil and the status is "success"
		return resp.Result == nil && resp.Status == "success"
	}

	// Execute the test with custom validation
	execute(t, retreiveUsers, validateFunc)
}

func retreiveUsersTestSuccess(t *testing.T) {
	// Define the expected validation function
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
			return false
		}

		// Check if the data is not nil and the status is "success"
		return resp.Result != nil && resp.Status == "success"
	}

	// Execute the test with custom validation
	execute(t, retreiveUsers, validateFunc)
}

func // CREATE
createUser(createData interface{}) func() (string, error) {
	return func() (string, error) {
		return action(
			"POST",
			endpoint+routeApiUsers,
			createData,
		)
	}
}
func // CREATE FAIL
createUserTestDuplicateDataFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, createUser(failUserCreateData), validateFunc)
}
func createUserTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, createUser(failUserCreateData), validateFunc)
}
func // CREATE SUCCESS
createUserTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, createUser(successUserCreateData), validateFunc)
}
func createUserTestSecondSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, createUser(successUserCreateSecondData), validateFunc)
}
func // RETRIEVE
retrieveUser() (string, error) {
	return action(
		"GET",
		endpoint+routeApiUsers+ID,
		nil,
	)
}
func // RETRIEVE SUCCESS
retrieveUserTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, retrieveUser, validateFunc)
}
func // RETRIEVE FAIL
retrieveUserTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, retrieveUser, validateFunc)
}
func // UPDATE
updateUser(updateData interface{}) func() (string, error) {
	return func() (string, error) {
		return action(
			"PUT",
			endpoint+routeApiUsers+ID,
			updateData,
		)
	}
}
func // UPDATE FAIL
updateUserTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, updateUser(failUserUpdateData), validateFunc)
}
func // UPDATE SUCCESS
updateUserTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, updateUser(successUserUpdateData), validateFunc)
}
func // DELETE
deleteUser() (string, error) {
	return action(
		"DELETE",
		endpoint+routeApiUsers+ID,
		nil,
	)
}
func // DELETE FAIL
deleteUserTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "error"
	}
	execute(t, deleteUser, validateFunc)
}
func // DELETE SUCCESS
deleteUserTestSuccess(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, deleteUser, validateFunc)
}

var // DATA
(
	// Item Test Data
	//
	// Fail cases
	// we don't store empty content
	failItemCreateData = models.JSONItemData{
		Title:       "title",
		Description: "",
	}
	// we don't store empty titles
	failItemUpdateData = models.JSONItemData{
		Title:       "",
		Description: "Data",
	}

	// Success cases
	successItemCreateData = models.JSONItemData{
		Title:       "First Post",
		Description: "love you Joyce",
		Image:       img,
	}
	successItemUpdateData = models.JSONItemData{
		Title:       "squirrel",
		Description: "Some words to drive you nuts",
		Image:       "",
	}

	// Fail cases
	// resopnse, err := dero.GetEncryptedBalance(address)
	// response.Result.Status != "OK"
	failCreateAddress  = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj"
	failUpdateAddress  = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0"
	failUserCreateData = models.User{
		Name:   user,
		Wallet: failCreateAddress,
	}
	failUserUpdateData = models.User{
		Name:   user,
		Wallet: failUpdateAddress,
	}
	// Success cases
	successCreateAddress       = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj8"
	successCreateSecondAddress = "dero1qykz2fqtptcnvmr65042jwljpwmglnezax4wms5w4htat20vzsdauqq58979y"
	successUpdateAddress       = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"

	successUserCreateData = models.User{
		Name:   user,
		Wallet: successCreateAddress,
	}
	successUserCreateSecondData = models.User{
		Name:   user2,
		Wallet: successCreateSecondAddress,
	}
	successUserUpdateData = models.User{
		Name:   user,
		Wallet: successUpdateAddress,
	}
)

var ( // base64encoded image
	img = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
)
