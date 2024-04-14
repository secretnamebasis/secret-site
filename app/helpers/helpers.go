package helpers

import "github.com/secretnamebasis/secret-site/app/database"

// Define bucket names
const (
	bucketItems = "items"
	bucketUsers = "users"
)

// NextItemID returns the next available item ID.
func NextItemID() (int, error) {
	return database.NextID(bucketItems)
}

// NextUserID returns the next available user ID.
func NextUserID() (int, error) {
	return database.NextID(bucketUsers)
}
