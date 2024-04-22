package controllers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// Define bucket names
const (
	bucketItems = "items"
	bucketUsers = "users"
)

// CreateItemRecord creates a new item in the database.
func CreateItemRecord(order *models.JSON_Item_Order) (models.Item, error) {
	order.Validate()
	// Let's create an item
	var item models.Item

	// Validate the order data
	if err := order.Validate(); err != nil {
		return models.Item{}, err
	}
	item.Title = order.Title
	// we are going to check by title...
	// no duplicate titles allowed
	if err := checkItemExistence(item.Title); err != nil {
		return models.Item{}, err
	}

	// Get the next item ID
	id, err := NextItemID()
	if err != nil {
		return models.Item{}, err
	}
	item.ID = id

	// Marshal the JSON_Item_Order into bytes
	// this is a really important concept:
	// we are going to be doing and seeing this kind
	// of operation a lot, we are going to be
	// marshalling data into bytes into some kind
	// of model so we are taking an order and we
	// are effectively building the ItemData that
	// will be stored as bytes.

	// and our validation already checks to see if
	// these fields are empty
	bytes, err := json.Marshal(
		models.ItemData{
			Description: order.Description,
			Image:       order.Image,
		},
	)
	// this should not be a problem...
	// but if it is...
	if err != nil {
		return models.Item{}, err
	}

	// Ensure that the marshaled bytes are not nil
	if bytes == nil {
		return models.Item{}, errors.New("marshaled bytes are nil")
	}

	// strap those bytes to our item
	item.Data = bytes
	// Encrypt bytes before storing in the database
	encryptedBytes,
		err := cryptography.EncryptData(
		item.Data, // we want to lock these bitches down!
		config.Env( // and to do it we are going into our env
			"SECRET", // and we are going to refer to our secret
		), // shouldn't ever change unless we are changing our encryption scheme
	)
	// and it better work.
	if err != nil {
		return models.Item{}, err
	}
	item.Data = encryptedBytes
	item.Initialize()

	// Validate the item
	if err := item.Validate(); err != nil {
		return models.Item{}, err
	}

	// Create the item record in the database
	err = database.CreateRecord(bucketItems, &item)
	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

// CreateUserRecord creates a new user in the database.
func CreateUserRecord(order *models.JSON_User_Order) error {

	// validate the wallet before moving forward
	if err := isValidWallet(order.Wallet); err != nil {
		return err
	}

	// we can't validate for existence in the model because of
	// a restriction on import cycle:
	// models <- database <- models
	// Controller check if user already exists with the same username or wallet
	if err := checkUserExistence(*order); err != nil {
		return err
	}

	var user models.User

	timestamp := time.Now()
	// Get the next item ID
	id, err := NextUserID()
	if err != nil {
		return err
	}
	user.ID = id
	user.Name = order.Name
	user.Wallet = order.Wallet
	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp
	// we store passwords as encrypted bytes for now
	// it would be better to hash their pass
	// and then compare the encrypted hash we have on file
	// against the pasword that they give us as hash
	encryptedPassword, err := cryptography.EncryptData(
		[]byte(order.Password),
		config.Env("SECRET"),
	)
	if err != nil {
		return err
	}
	user.Password = encryptedPassword

	// Validate wallet address
	if err := user.Validate(); err != nil {
		return err
	}
	// Create the user with the provided data
	user.Initialize()

	// Store the user record in the database
	return database.CreateRecord(bucketUsers, &user)
}

// isValidWallet checks if the provided wallet address is valid
func isValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(config.NodeEndpoint, wallet)
	return err
}

// AllItems retrieves all items from the database.
func AllItems() ([]models.Item, error) {
	var items []models.Item
	err := database.GetAllRecords(bucketItems, &items)
	if err != nil {
		return nil, err // Return nil slice and error
	}
	return items, nil // Return retrieved items and nil error
}

// AllItems retrieves all items from the database.
func AllItemTitles() ([]models.Item, error) {
	var items []models.Item
	err := database.GetAllItemTitles(bucketItems, &items)
	if err != nil {
		return nil, err // Return nil slice and error
	}
	return items, nil // Return retrieved items and nil error
}

// AllUsers retrieves all users from the database.
func AllUsers() ([]models.User, error) {
	var users []models.User
	err := database.GetAllRecords(bucketUsers, &users)
	return users, err
}

// GetUserByID retrieves a user from the database by ID.
func GetUserByID(id string) (models.User, error) {
	var user models.User
	err := database.GetRecordByID(bucketUsers, id, &user)
	return user, err
}

func GetUserByName(name string) (models.User, error) {
	existingUser, err := database.GetUserByUsername(name)
	if err != nil {
		return *existingUser, errors.New("error checking user existence")
	}
	if existingUser != nil {
		return *existingUser, nil
	}
	return models.User{}, err
}

// GetItemByID retrieves an item from the database by ID.
func GetItemByID(id string) (models.Item, error) {
	var item models.Item
	err := database.GetRecordByID(bucketItems, id, &item)
	return item, err
}

// UpdateItem updates an item in the database with the provided ID and updated data.
func UpdateItem(id string, updatedItem models.JSON_Item_Order) error {
	return database.UpdateRecord(bucketItems, id, &updatedItem)
}

// UpdateUser updates a user in the database with the provided ID and updated data.
func UpdateUser(id string, updatedUser models.User) error {
	// Check if user with the provided ID exists
	existingUser, err := database.GetUserByUsername(updatedUser.Name)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Validate wallet address
	if err := validateWalletAddress(updatedUser.Wallet); err != nil {
		return err
	}

	// Update the user record in the database
	return database.UpdateRecord(bucketUsers, id, &updatedUser)
}

// DeleteItem deletes an item from the database by ID.
func DeleteItem(id string) error {
	return database.DeleteRecord(bucketItems, id)
}

// DeleteUser deletes a user from the database by ID.
func DeleteUser(id string) error {
	return database.DeleteRecord(bucketUsers, id)
}

// NextUserID returns the next available user ID.
func NextUserID() (int, error) {
	return database.NextID(bucketUsers)
}

// NextItemID returns the next available item ID.
func NextItemID() (int, error) {
	return database.NextID(bucketItems)
}

// private functions
// checkItemExistence checks if a user with the same title or data already exists
func checkItemExistence(title string) error {

	// Check if user already exists with the same username
	existingItem, err := database.GetItemByTitle(title)
	if err != nil {
		return errors.New("error checking item existence")
	}

	if existingItem != nil {
		return errors.New("item with the same title already exists")
	}

	return nil
}

// checkUserExistence checks if a user with the same username or wallet already exists
func checkUserExistence(order models.JSON_User_Order) error {

	// Check if user already exists with the same username
	existingUser, err := database.GetUserByUsername(order.Name)
	if err != nil {
		return errors.New("error checking user existence")
	}

	if existingUser != nil {
		return errors.New("user with the same username already exists")
	}

	// Check if user already exists with the same wallet
	existingUser, err = database.GetUserByWallet(order.Wallet)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser != nil {
		return errors.New("user with the same wallet already exists")
	}

	return nil
}

// validateWalletAddress checks if the provided wallet address is valid
func validateWalletAddress(wallet string) error {
	if err := isValidWallet(wallet); err != nil {
		return errors.New("invalid wallet address")
	}
	return nil
}
