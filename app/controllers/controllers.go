package controllers

import (
	"encoding/json"
	"errors"

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
func CreateItemRecord(call *models.JSONItemData) (models.Item, error) {
	if err := checkItemExistence(call); err != nil {
		return models.Item{}, err
	}
	// Validate the input data
	if err := call.Validate(); err != nil {
		return models.Item{}, err
	}

	// Marshal the JSONItemData into bytes
	bytes, err := json.Marshal(models.ItemData{
		Description: call.Description,
		Image:       call.Image,
	})
	if err != nil {
		return models.Item{}, err
	}

	// Ensure that the marshaled bytes are not nil
	if bytes == nil {
		return models.Item{}, errors.New("marshaled bytes are nil")
	}

	// Create a new item
	var item models.Item
	// Get the next item ID
	id, err := NextItemID()
	if err != nil {
		return models.Item{}, err
	}
	item.ID = id
	item.Title = call.Title
	item.Data = bytes
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
func CreateUserRecord(user *models.User) error {

	// we can't validate for existence in the model because of
	// a restriction on import cycle:
	// models <- database <- models
	// Controller check if user already exists with the same username or wallet
	if err := checkUserExistence(*user); err != nil {
		return err
	}

	id, err := NextUserID()
	if err != nil {
		return err
	}
	user.ID = id

	// Validate wallet address
	if err := user.Validate(); err != nil {
		return err
	}
	// Create the user with the provided data
	user.Initialize()

	// Store the user record in the database
	return database.CreateRecord(bucketUsers, user)
}

// isValidWallet checks if the provided wallet address is valid
func isValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(wallet)
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

// GetItemByID retrieves an item from the database by ID.
func GetItemByID(id string) (models.Item, error) {
	var item models.Item
	err := database.GetRecordByID(bucketItems, id, &item)
	return item, err
}

// UpdateItem updates an item in the database with the provided ID and updated data.
func UpdateItem(id string, updatedItem models.JSONItemData) error {
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
func checkItemExistence(item *models.JSONItemData) error {

	// Check if user already exists with the same username
	existingItem, err := database.GetItemByTitle(item.Title)
	if err != nil {
		return errors.New("error checking item existence")
	}

	if existingItem != nil {
		return errors.New("item with the same title already exists")
	}

	return nil
}

// checkUserExistence checks if a user with the same username or wallet already exists
func checkUserExistence(user models.User) error {

	// Check if user already exists with the same username
	existingUser, err := database.GetUserByUsername(user.Name)
	if err != nil {
		return errors.New("error checking user existence")
	}

	if existingUser != nil {
		return errors.New("user with the same username already exists")
	}

	// Check if user already exists with the same wallet
	existingUser, err = database.GetUserByWallet(user.Wallet)
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
