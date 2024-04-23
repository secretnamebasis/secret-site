package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUserRecord creates a new user in the database.
func CreateUserRecord(order *models.JSON_User_Order) error {
	order.Validate()

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
	fmt.Printf("NEXT ID %d\n", id)
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
		config.Env(
			config.EnvPath,
			"SECRET",
		),
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

// UpdateUser updates a user in the database with the provided ID and updated data.
func UpdateUser(updatedUser models.User) error {
	// Check if user with the provided ID exists
	existingUser, err := database.GetUserByUsername(updatedUser.Name)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Validate wallet address
	if err := ValidateWalletAddress(updatedUser.Wallet); err != nil {
		return err
	}
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
	}
	if updatedUser.Wallet != "" {
		existingUser.Wallet = updatedUser.Wallet
	}
	// Always update the password if provided
	if updatedUser.Password != nil {
		existingUser.Password = updatedUser.Password
	} else {
		// If password is not provided, preserve the existing password
		updatedUser.Password = existingUser.Password
	}
	// Preserve the ID and creation timestamp
	updatedUser.ID = existingUser.ID
	updatedUser.CreatedAt = existingUser.CreatedAt
	updatedUser.UpdatedAt = time.Now()
	// Update the user record in the database
	return database.CreateRecord(bucketUsers, &updatedUser)
}

// DeleteUser deletes a user from the database by ID.
func DeleteUser(id string) error {
	return database.DeleteRecord(bucketUsers, id)
}

// NextUserID returns the next available user ID.
func NextUserID() (int, error) {
	return database.NextID(bucketUsers)
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
