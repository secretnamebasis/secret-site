package controllers

import (
	"errors"
	"time"

	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateUserRecord creates a new user in the database.
func CreateUserRecord(order *models.JSON_User_Order) error {
	// fmt.Println("order receieved")
	order.Validate()
	// fmt.Println("order validated")
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

	user = models.User{
		ID:        id,
		Name:      order.Name,
		Wallet:    order.Wallet,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	// we are no longer doing this because the tx is the password
	//
	// user.Password = cryptography.HashString( // so let's hash the string up
	// 	order.Password, // because we don't want to record this anywhere
	// )

	// Validate wallet address
	if err := user.Validate(); err != nil {
		return err
	}
	// fmt.Println("user validated")

	// Create the user with the provided data
	user.Initialize()
	// fmt.Println("user initialized")
	// fmt.Println(user)
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

// GetUserByID retrieves a user from the database by ID.
func GetUserByWallet(wallet string) (models.User, error) {
	user, err := database.GetUserByWallet(wallet)
	return user, err
}

func GetUserByName(name string) (models.User, error) {
	existingUser, err := database.GetUserByUsername(name)
	if err != nil {
		return existingUser, errors.New("error checking user existence")
	}
	if existingUser.Name != "" {
		return existingUser, nil
	}
	return models.User{}, err
}

// UpdateUser updates a user in the database with the provided ID and updated data.
func UpdateUser(order models.JSON_User_Order) error {
	// Check if user with the provided ID exists
	existingUser, err := database.GetUserByUsername(order.Name)

	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser.Name == "" {
		return errors.New("user not found")
	}

	// Validate wallet address
	if err := ValidateWalletAddress(order.Wallet); err != nil {
		return err
	}
	if order.Name != "" {
		existingUser.Name = order.Name
	}
	if order.Wallet != "" {
		existingUser.Wallet = order.Wallet
	}

	// we are removing passwords
	// if order.Password != "" {
	// 	existingUser.Password = []byte( // it will be "best" to store as byte
	// 		cryptography.HashString( // so let's hash the string up
	// 			order.Password, // because we don't want to record this anywhere
	// 		),
	// 	)
	// }

	existingUser.UpdatedAt = time.Now()

	// Update the user record in the database
	return database.CreateRecord(bucketUsers, &existingUser)
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

	if existingUser.Name != "" {
		return errors.New("user with the same username already exists")
	}

	// Check if user already exists with the same wallet
	existingUser, err = database.GetUserByWallet(order.Wallet)
	if err != nil {
		return errors.New("error checking user existence")
	}
	if existingUser.Name != "" {
		return errors.New("user with the same wallet already exists")
	}

	return nil
}
