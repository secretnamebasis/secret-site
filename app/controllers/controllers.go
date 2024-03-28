package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/db"
	"github.com/secretnamebasis/secret-site/app/models"
)

// Define bucket names
const (
	bucketItems = "items"
	bucketUsers = "users"
)

// CreateItem creates a new item in the database.
func CreateItem(item models.Item) error {
	return db.CreateRecord(bucketItems, item)
}

// CreateUser creates a new user in the database.
func CreateUser(user models.User) error {
	return db.CreateRecord(bucketUsers, user)
}

// AllItems retrieves all items from the database.
func AllItems(c *fiber.Ctx) ([]models.Item, error) {
	var items []models.Item
	err := db.GetAllRecords(bucketItems, &items, c)
	return items, err
}

// AllUsers retrieves all users from the database.
func AllUsers(c *fiber.Ctx) ([]models.User, error) {
	var users []models.User
	err := db.GetAllRecords(bucketUsers, &users, c)
	return users, err
}

// GetUserByID retrieves a user from the database by ID.
func GetUserByID(id string) (models.User, error) {
	var user models.User
	err := db.GetRecordByID(bucketUsers, id, &user)
	return user, err
}

// GetItemByID retrieves an item from the database by ID.
func GetItemByID(id string) (models.Item, error) {
	var item models.Item
	err := db.GetRecordByID(bucketItems, id, &item)
	return item, err
}

// UpdateItem updates an item in the database with the provided ID and updated data.
func UpdateItem(id string, updatedItem models.Item) error {
	return db.UpdateRecord(bucketItems, id, &updatedItem)
}

// UpdateUser updates a user in the database with the provided ID and updated data.
func UpdateUser(id string, updatedUser models.User) error {
	return db.UpdateRecord(bucketUsers, id, &updatedUser)
}

// DeleteItem deletes an item from the database by ID.
func DeleteItem(id string) error {
	return db.DeleteRecord(bucketItems, id)
}

// DeleteUser deletes a user from the database by ID.
func DeleteUser(id string) error {
	return db.DeleteRecord(bucketUsers, id)
}

// NextUserID returns the next available user ID.
func NextUserID() (int, error) {
	return db.NextID(bucketUsers)
}

// NextItemID returns the next available item ID.
func NextItemID() (int, error) {
	return db.NextID(bucketItems)
}
