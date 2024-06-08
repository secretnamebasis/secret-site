package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

// CreateItemRecord creates a new item in the database.
func CreateItemRecord(order *models.JSON_Item_Order) (models.Item, error) {

	// if err := authenticateUser(order.User); err != nil {
	// 	return models.Item{}, err
	// }

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

	if _, err := dero.GetSCID(config.NodeEndpoint, order.SCID); err != nil {
		return models.Item{}, err
	}
	item.SCID = order.SCID

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
			File:        order.File,
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
		config.Env(
			config.EnvPath, // and to do it we are going into our env
			"SECRET",       // and we are going to refer to our secret
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

// GetItemByID retrieves an item from the database by ID.
func GetItemByID(id string) (models.Item, error) {
	var existingItem models.Item

	if err := database.GetRecordByID(bucketItems, id, &existingItem); err != nil {
		return models.Item{}, err
	}

	decryptedData, // seeing as this is a big garbaldy goop...
		err := cryptography.DecryptData(
		existingItem.Data,
		config.Env(config.EnvPath, "SECRET"),
	)
	if err != nil {
		return models.Item{}, err
	}
	existingItem.Data = decryptedData
	return existingItem, err
}

// GetItemBySCID retrieves an item from the database by SCID.
func GetItemBySCID(scid string) (models.Item, error) {

	item, err := database.GetItemByField("scid", scid)
	if err != nil {
		return models.Item{}, err
	}

	decryptedData, // seeing as this is a big garbaldy goop...
		err := cryptography.DecryptData(
		item.Data,
		config.Env(config.EnvPath, "SECRET"),
	)
	if err != nil {
		return models.Item{}, err
	}
	item.Data = decryptedData
	return item, err
}

// UpdateItem updates an item in the database with the provided ID and updated data.
func UpdateItem(scid string, order models.JSON_Item_Order) error {

	// you are trying to get rid of passwords
	// if err := authenticateUser(order.User); err != nil {
	// 	return err
	// }

	existingItem, err := GetItemBySCID(order.SCID)
	if err != nil {
		return err
	}
	// if err := database.GetRecordByID(bucketItems, id, &existingItem); err != nil {
	// 	return err
	// }

	// decryptedData, // seeing as this is a big garbaldy goop...
	// 	err := cryptography.DecryptData(
	// 	existingItem.Data,
	// 	config.Env(
	// 		config.EnvPath,
	// 		"SECRET",
	// 	),
	// )
	// if err != nil {
	// 	return err
	// }
	// let's go put this all back together
	var existingItemData models.ItemData
	if err := json.Unmarshal(
		existingItem.Data,
		&existingItemData,
	); err != nil {
		return err
	}
	if order.Title != "" {
		existingItem.Title = order.Title
	}
	// Update the existingItemData fields
	if order.Image != "" {
		existingItemData.Image = order.Image
	}
	if order.Description != "" {
		existingItemData.Description = order.Description
	}

	if order.File != "" {
		existingItemData.File = order.File
	}
	if order.Image != "" {
		existingItemData.Image = order.Image
	}

	// Marshal the updated data and encrypt it
	updatedBytes, err := json.Marshal(existingItemData)
	if err != nil {
		return err
	}

	// currently we are encrypting with our password
	// we could futher hash the password for greater security...
	// or we could use DERO
	encryptedBytes,
		err := cryptography.EncryptData(
		updatedBytes,
		config.Env(
			config.EnvPath,
			"SECRET",
		),
	)
	if err != nil {
		return err
	}

	// Update existingItem with the encrypted data and set the updated timestamp
	existingItem.Data = encryptedBytes
	existingItem.UpdatedAt = time.Now()

	return database.CreateRecord(bucketItems, &existingItem)
}

// DeleteItem deletes an item from the database by ID.
func DeleteItem(scid string) error {
	return database.DeleteRecord(bucketItems, scid)
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

	if existingItem.Title != "" {
		return errors.New("item with the same title already exists")
	}

	if existingItem.SCID != "" {
		return errors.New("item with the same scid already exists")
	}

	return nil
}

// authenticateUser checks if a user with the same username or wallet already exists
func authenticateUser(order models.JSON_User_Order) error {

	// Check if a user already exists with the same username
	existingUser, err := database.GetUserByUsername(order.Name)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return errors.New("error checking user existence")
	}
	if existingUser.Name == "" {
		log.Printf("user does not exist: %v", err)
		return errors.New("user does not exist")
	}
	// we are not using passwords anymore
	// // Hash the password for comparison
	// hashedPass := cryptography.HashString(
	// 	order.Password,
	// )

	// // Compare hashed passwords to authenticate
	// if !bytes.Equal(existingUser.Password, hashedPass) {
	// 	log.Println("Invalid password")
	// 	return errors.New("error invalid password")
	// }

	return nil
}
