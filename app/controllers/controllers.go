package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// Define bucket names
const (
	bucketItems = "items"
	bucketUsers = "users"
)

// CreateItem creates a new item in the database.
func CreateItemRecord(item *models.Item) error {
	return database.CreateRecord(bucketItems, item)
}

// CreateUserRecord creates a new user in the database.
func CreateUserRecord(user *models.User) error {
	// Check if user already exists with the same username or wallet
	if err := checkUserExistence(*user); err != nil {
		return err
	}

	// Validate wallet address
	if err := user.Validate(); err != nil {
		return err
	}

	// Generate ID and password
	nextID, _ := NextUserID()
	password := uuid.New().String()

	// Create the user with the provided data
	newUser := &models.User{
		ID:        nextID,
		User:      user.User,
		Wallet:    user.Wallet,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Store the user record in the database
	return database.CreateRecord(bucketUsers, newUser)
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
func AllUsers(c *fiber.Ctx) ([]models.User, error) {
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
func UpdateItem(id string, updatedItem models.Item) error {
	return database.UpdateRecord(bucketItems, id, &updatedItem)
}

// UpdateUser updates a user in the database with the provided ID and updated data.
func UpdateUser(id string, updatedUser models.User) error {
	// Check if user with the provided ID exists
	existingUser, err := database.GetUserByUsername(updatedUser.User)
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

// checkUserExistence checks if a user with the same username or wallet already exists
func checkUserExistence(user models.User) error {

	// Check if user already exists with the same username
	existingUser, err := database.GetUserByUsername(user.User)
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
