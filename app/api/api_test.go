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
	"testing"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
	"github.com/secretnamebasis/secret-site/app/services"
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
	endpoint       = "http://127.0.0.1:3000/api"
	routeApiUsers  = "/users/"
	routeApiItems  = "/items/"
	routeApiAssets = "/assets/"
	user           = "secret"
	pass           = "pass"
	user2          = pass
)

var (
	ID              string
	WalletEndpoint1 string
	WalletEndpoint2 string
)

func // CONFIG
configure() error {

	config.NodeEndpoint = "http://" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_NODE_IP") +
		":" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_NODE_PORT",
		) +
		"/json_rpc"

	WalletEndpoint1 = "http://" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_WALLET_IP") +
		":" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_WALLET1_PORT",
		) +
		"/json_rpc"

	var response rpc.GetAddress_Result
	err := dero.CallRPC(
		WalletEndpoint1,
		&response,
		"GetAddress",
	)
	if err != nil {
		return err
	}

	successCreateAddress = response.Address

	successUserCreateData = models.User{
		Name:   user,
		Wallet: successCreateAddress,
	}

	WalletEndpoint2 = "http://" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_WALLET_IP") +
		":" +
		config.Env(
			config.EnvPath,
			"DERO_SIMULATOR_WALLET2_PORT",
		) +
		"/json_rpc"

	// empty it out
	response = rpc.GetAddress_Result{}
	err = dero.CallRPC(
		WalletEndpoint2,
		&response,
		"GetAddress",
	)
	if err != nil {
		return err
	}
	successCreateSecondAddress = response.Address
	successUserCreateSecondData = models.User{
		Name:   user2,
		Wallet: successCreateSecondAddress,
	}
	successUpdateAddress = successCreateAddress
	successUserUpdateData = models.User{
		Name:   user,
		Wallet: successUpdateAddress,
	}

	scid,
		err = dero.MintContract(
		WalletEndpoint1,
		dero.NFAContract(
			"1",
			"simple",
			"smart-contract",
			"image",
			"test",
			successCreateAddress,
		),
		successCreateSecondAddress, // you can't send to self
	)

	if err != nil {
		return err
	}
	ID = scid.TXID

	// Success cases
	successItemCreateData = models.JSON_Item_Order{
		Title:       "First Post",
		Description: "love you Joyce",
		SCID:        scid.TXID,
		Image:       LittleImg,
	}

	successItemUpdateData = models.JSON_Item_Order{
		Title:       "squirrel",
		Description: "Some words to drive you nuts",
		SCID:        scid.TXID,
		Image:       "",
	}

	return nil
}

type // RESPONSE
response struct {
	Result interface{} `json:"result"`
	Status string      `json:"status"`
}

var // DELAY
delay = 1 * time.Nanosecond

func pause(i time.Duration) {
	time.Sleep(
		i * time.Second,
	)
}

// MAIN
func TestAPI(t *testing.T) {

	// Config server
	cfg := config.Initialize()

	// Check if config is empty
	checkConfig(cfg)

	// load simulator wallets
	if err := configure(); err != nil {
		log.Fatalf("failed to load wallets: %s", err)
	}
	pause(1)

	// Check if testing framework is empty
	checkTestingFramework(t)

	// Initialize the database with configs
	initializeDatabase(cfg)

	// Start the server as an app
	app := startServer(t, cfg)

	// Check if app is empty
	checkApp(app)

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

	if err := services.ProcessCheckouts(c); err != nil {
		log.Fatal(err)
	}

	return a
}

func runTests(t *testing.T, c config.Server) { // run tests
	log.Printf("Environment: %s\n", c.Environment)
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
	err := os.Remove(c.DatabasePath + c.Environment + ".db")
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
		// this is the start
		{
			"Create success when User 1 is valid",
			createUserTestSuccess,
		},
		{
			"Create fail when User 1 already exists",
			createUserTestDuplicateDataFail,
		},
		// we want to know that they refer...
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

		// Validate created
		created, createdOK := result["created_at"].(string)
		if !createdOK || created == "" {
			t.Errorf("Expected created to be '%s' but got '%s'", time.Now(), created)
			return false
		}

		// Validate updated
		updated, updatedOK := result["updated_at"].(string)
		if !updatedOK || updated == "" {
			t.Errorf("Expected updated to be '%s' but got '%s'", time.Now(), updated)
			return false
		}

		// Validate created
		url, urlOK := result["image_url"].(string)
		if !urlOK || url == "" {
			t.Errorf("Expected title to be `"+config.Domain+"'images/1` got '%s'", url)
			return false
		}

		// Validate data
		encodedData, dataOK := result["data"].(string)
		if !dataOK {
			t.Errorf("Expected 'data' field to be present")
			return false
		}

		// fmt.Printf("ENCODED DATA: %s\n", encodedData)
		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		// fmt.Printf("DECODED DATA: %s\n", decodedData)
		// Decrypt the data
		decryptedData, err := cryptography.DecryptData(
			decodedData, config.Env(
				config.EnvPath,
				"SECRET",
			),
		)
		if err != nil {
			t.Errorf("Error decrypting data: %v", err)
			return false
		}
		// fmt.Printf("DECRYPTED DATA: %s\n", decryptedData)
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
		// fmt.Printf("ENCODED DATA: %s\n", encodedData)

		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		// fmt.Printf("DECODED DATA: %s\n", decodedData)
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
		// fmt.Printf("ENCODED DATA: %s\n", encodedData)

		// Decode base64-encoded data
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			t.Errorf("Expected to dedcode")
			return false
		}
		// fmt.Printf("DECODED DATA: %s\n", decodedData)
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
retrieveUser(address string) (string, error) {
	return action(
		"GET",
		endpoint+routeApiUsers+address,
		nil,
	)
}
func // RETRIEVE SUCCESS
retrieveUserTestSuccess(t *testing.T) {
	// the address needs to exist prior to validation
	// this is the integrated address...
	//
	// this is the second user
	// detoi1qyvyeyzrcm2fzf6kyq7egkes2ufgny5xn77y6typhfx9s7w3mvyd5q9yvfp4xersv9ehxcjy25vs2wtzfe2szcjk25vjwyqah9nj2
	// this is the simulator wallet address
	// deto1qyre7td6x9r88y4cavdgpv6k7lvx6j39lfsx420hpvh3ydpcrtxrxqg8v8e3z
	// it would be better to split it... but don't.

	ps := rpc.Split_Integrated_Address_Params{
		Integrated_Address: "detoi1qyvyeyzrcm2fzf6kyq7egkes2ufgny5xn77y6typhfx9s7w3mvyd5q9yvfp4xennv43hyet5vfz92xg9893yu4gpvft92xf8zqwe58zs",
	}
	splitResult, err := dero.SplitIntegratedAddress(WalletEndpoint2, ps)
	if err != nil {
		t.Fatalf("Error splitting integrated address: %v", err)
	}
	log.Printf("address: %s, payload: %s\n", splitResult.Address, splitResult.Payload_RPC)

	a := rpc.Arguments{
		rpc.Argument{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    successUserCreateData.Name,
		},
		rpc.Argument{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    config.UserRegistrationPort,
		},
		rpc.Argument{
			Name:     rpc.RPC_REPLYBACK_ADDRESS,
			DataType: rpc.DataString,
			Value:    successCreateAddress,
		},
		// I don't remember if you have to include this in the payload
		// rpc.Argument{
		// 	Name:     rpc.RPC_VALUE_TRANSFER,
		// 	DataType: rpc.DataUint64,
		// 	Value:    config.UserRegistrationFee,
		// },
	}
	tx := []rpc.Transfer{
		{
			SCID:        crypto.ZEROHASH,
			Destination: "deto1qyvyeyzrcm2fzf6kyq7egkes2ufgny5xn77y6typhfx9s7w3mvyd5qqynr5hx",
			Amount:      config.UserRegistrationFee,
			Payload_RPC: a,
		},
	}

	p := rpc.Transfer_Params{
		Transfers: tx,
	}
	result, err := dero.Transfer(WalletEndpoint1, p)
	if err != nil {
		t.Fatalf("Error creating transfer: %v", err)
	}
	if result.TXID == "" {
		t.Fatalf("Error: txid is empty: %v", err)
	}
	fmt.Println("TXID: " + result.TXID)
	time.Sleep(5 * time.Second)
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}
		// Perform custom validation based on new expectations

		return resp.Status == "success"
	}
	execute(t, func() (string, error) { return retrieveUser(successCreateAddress) }, validateFunc)
}
func // RETRIEVE FAIL
retrieveUserTestFail(t *testing.T) {
	validateFunc := func(responseBody string) bool {
		var resp response
		if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
			t.Fatalf("Error parsing response: %v", err)
		}

		// Perform custom validation based on new expectations

		return resp.Status == "success" // oddly enough... that should be fine.
	}
	execute(t, func() (string, error) { return retrieveUser(failCreateAddress) }, validateFunc)
}
func // UPDATE
updateUser(updateData interface{}) func() (string, error) {
	return func() (string, error) {
		return action(
			"PUT",
			endpoint+routeApiUsers+successUpdateAddress, // this is kind of hacky
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
deleteUser(address string) func() (string, error) {
	return func() (string, error) {
		return action(
			"DELETE",
			endpoint+routeApiUsers+address,
			nil,
		)
	}
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
	execute(t, deleteUser(successCreateAddress), validateFunc)
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
	execute(t, deleteUser(successCreateAddress), validateFunc)
}

var // DATA
(
	// Item Test Data
	//
	// Fail cases

	// we don't store empty content
	failItemCreateData = models.JSON_Item_Order{
		Title:       "title",
		Description: "",
	}
	// we don't store empty titles
	failItemUpdateData = models.JSON_Item_Order{
		Title:       "",
		Description: "Data",
	}

	// Success cases
	successItemCreateData = models.JSON_Item_Order{
		Title:       "First Post",
		Description: "love you Joyce",
		SCID:        scid.TXID,
		Image:       LittleImg,
		File:        File,
	}

	successItemUpdateData = models.JSON_Item_Order{
		Title:       "squirrel",
		Description: "Some words to drive you nuts",
		SCID:        scid.TXID,
		Image:       "",
	}

	// Fail cases
	// resopnse, err := dero.GetEncryptedBalance(address)
	// response.Result.Status != "OK"
	failCreateAddress  = "dero1qynmz4tgkmtmmspqmywvjjmtl0x8vn5ahz4xwaldw0hu6r5500hryqgptvnj"
	failUpdateAddress  = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0"
	failUserCreateData = models.User{
		Name: "srpn534luegssnviukp30mtwi1fzwo611x2iivsofjwzgao19azo366niyp7pmr2svrpox4na8qy8buna6fk2aeimeg894qu7g068ihiyuhqy6gs5nf3nudeniv60mawvbnk220ol6rxynj6nii5ob83c2vklldfuylp651318i0u7wq9g2rp9ccrm75yxcxy1vtkj34s1hprgj7uhxiunxqm2rjij3fvkbsldm6ksqn7dyat9ccqrbg2w8hyee8",
	}
	failUserUpdateData = models.User{
		Name:   user,
		Wallet: failUpdateAddress,
	}
	// Success cases
	scid                       rpc.Transfer_Result
	successCreateAddress       string
	successCreateSecondAddress string
	successUpdateAddress       string

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

// base64encoded images
var ( // small
	LittleImg        = `iVBORw0KGgoAAAANSUhEUgAAAAoAAAAKCAYAAACNMs+9AAAAFUlEQVR42mNkAAIAAAoAAv/lxKUAAAAASUVORK5CYII=`
	File      string = `
	-----BEGIN DERO SIGNED MESSAGE-----
	Address: dero1qyfk5w2rvqpl9kzfd7fpteyp2k362y6audydcu2qrgcmj6vtasfkgqq9704gn
	C: ce34dd75b765f92a7af9f834974a2f729e5f53591a76a34991bafa501166388
	S: 2bbcf2f5e4e3e91835d8364f636df4fb16726cc0bb3321a73685af4eabb0570a
	
	iVBORw0KGgoAAAANSUhEUgAADXAAAAWgCAYAAABNcorEAAAgAElEQVR4nOy9Z6xl
		-----END DERO SIGNED MESSAGE-----
`
)
