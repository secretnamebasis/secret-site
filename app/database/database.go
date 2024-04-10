package database

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/secretnamebasis/secret-site/app/models"

	"go.etcd.io/bbolt"
)

var (
	db          *bbolt.DB
	itemsBucket = []byte("items")
	usersBucket = []byte("users")
)

func InitDB(c config.Server) error {
	var err error
	// set directory of the database
	if err := os.MkdirAll(c.DatabasePath, 0755); err != nil {
		return err
	}

	db, err = bbolt.Open(c.DatabasePath+c.Env+".db", 0600, nil)
	if err != nil {
		return err
	}

	// Create buckets if they don't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(itemsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(usersBucket)
		return err
	})

	return err
}

// CreateRecord creates a record in the database for the given item/user after encrypting the content.
func CreateRecord(bucketName string, record interface{}) error {
	return db.Update(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %q not found ", bucketName)
			}
			var id int
			var encryptedContent string
			switch r := record.(type) {
			case *models.Item: // Ensure we're dealing with a pointer to models.Item
				id = r.ID
				// Encrypt content before storing in the database
				encrypted, err := cryptography.EncryptData([]byte(r.Content.Description), config.Env(exports.EnvPath))
				if err != nil {
					return err
				}

				// Encode the encrypted content to Base64
				encryptedContent = base64.StdEncoding.EncodeToString(encrypted)
				r.Content.Description = encryptedContent

				encrypted, err = cryptography.EncryptData([]byte(r.Content.Image), config.Env(exports.EnvPath))
				if err != nil {
					return err
				}

				// Encode the encrypted content to Base64
				encryptedContent = base64.StdEncoding.EncodeToString(encrypted)
				r.Content.Image = encryptedContent

			case *models.User:
				id = r.ID

			default:
				return fmt.Errorf("unsupported record type")
			}

			recordJSON, err := json.Marshal(record)
			if err != nil {
				return err
			}
			return b.Put([]byte(strconv.Itoa(id)), recordJSON)
		},
	)
}

// GetAllRecords retrieves all records from the specified bucket and unmarshals them into the provided slice.
func GetAllRecords(bucketName string, records interface{}) error {
	return db.View(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %q not found", bucketName)
			}

			// Define a helper function to unmarshal records
			unmarshalRecord := func(recordType interface{}) error {
				return b.ForEach(
					func(k, v []byte) error {
						if err := json.Unmarshal(v, recordType); err != nil {
							return err
						}
						// Decrypt the content if the record type is models.Item
						if item, ok := recordType.(*models.Item); ok {
							decodedBytes, err := base64.StdEncoding.DecodeString(item.Content.Description)
							if err != nil {
								return err
							}
							decryptedContent, err := cryptography.DecryptData(decodedBytes, config.Env(exports.EnvPath))
							if err != nil {
								return err
							}
							item.Content.Description = string(decryptedContent)
							decodedBytes, err = base64.StdEncoding.DecodeString(item.Content.Image)
							if err != nil {
								return err
							}
							decryptedContent, err = cryptography.DecryptData(decodedBytes, config.Env(exports.EnvPath))
							if err != nil {
								return err
							}
							item.Content.Image = string(decryptedContent)
						}
						sliceValue := reflect.ValueOf(records).Elem()
						sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(recordType).Elem()))
						return nil
					})
			}

			// Unmarshal records based on their types
			switch records.(type) {
			case *[]models.Item:
				return unmarshalRecord(&models.Item{})
			case *[]models.User:
				return unmarshalRecord(&models.User{})
			default:
				return fmt.Errorf("unsupported record type")
			}
		})
}

// UpdateRecord updates a record in the specified bucket with the provided ID and updated data.
func UpdateRecord(bucketName, id string, updatedRecord interface{}) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		recordJSON := b.Get([]byte(id))
		if recordJSON == nil {
			return fmt.Errorf("record with ID %s not found in bucket %s", id, bucketName)
		}

		if err := updateRecordFromJSON(recordJSON, updatedRecord); err != nil {
			return err
		}

		updatedJSON, err := json.Marshal(updatedRecord)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), updatedJSON)
	})
}

// updateRecordFromJSON updates the provided record using data from JSON.
func updateRecordFromJSON(recordJSON []byte, updatedRecord interface{}) error {
	switch record := updatedRecord.(type) {
	case *models.Item:
		var existingItem models.Item
		if err := json.Unmarshal(recordJSON, &existingItem); err != nil {
			return err
		}
		updateExistingItem(record, &existingItem)
	case *models.User:
		var existingUser models.User
		if err := json.Unmarshal(recordJSON, &existingUser); err != nil {
			return err
		}
		updateExistingUser(record, &existingUser)
	default:
		return fmt.Errorf("unsupported record type")
	}
	return nil
}

// updateExistingItem updates the item record based on existing data.
func updateExistingItem(updatedItem, existingItem *models.Item) {
	// Define a map to store the fields to be updated
	fieldsToUpdate := map[string]bool{
		"Content.ImageURL":    updatedItem.Content.ImageURL != "",
		"Content.Image":       updatedItem.Content.Image != "",
		"Title":               updatedItem.Title != "",
		"Content.Description": updatedItem.Content.Description != "",
	}

	// Iterate over the fields of the updatedItem and update the corresponding fields in the existingItem
	updatedItemType := reflect.TypeOf(*updatedItem)
	updatedItemValue := reflect.ValueOf(*updatedItem)
	existingItemValue := reflect.ValueOf(existingItem).Elem() // Dereference the pointer to access the struct fields

	for i := 0; i < updatedItemType.NumField(); i++ {
		field := updatedItemType.Field(i)
		if fieldsToUpdate[field.Name] {
			fieldValue := updatedItemValue.Field(i)
			existingItemValue.FieldByName(field.Name).Set(fieldValue)
		}
	}

	// Preserve the ID, creation timestamp, and update timestamp
	existingItem.ID = updatedItem.ID
	existingItem.CreatedAt = existingItem.CreatedAt
	existingItem.UpdatedAt = time.Now()
}

// updateUserFromExisting updates the user record based on existing data.
func updateExistingUser(updatedUser, existingUser *models.User) {
	if updatedUser.User != "" {
		existingUser.User = updatedUser.User
	}
	if updatedUser.Wallet != "" {
		existingUser.Wallet = updatedUser.Wallet
	}
	// Always update the password if provided
	if updatedUser.Password != "" {
		existingUser.Password = updatedUser.Password
	} else {
		// If password is not provided, preserve the existing password
		updatedUser.Password = existingUser.Password
	}
	// Preserve the ID and creation timestamp
	updatedUser.ID = existingUser.ID
	updatedUser.CreatedAt = existingUser.CreatedAt
	updatedUser.UpdatedAt = time.Now()
}

// GetRecordByID retrieves a record from the specified bucket by ID and unmarshals it into the provided model.
func GetRecordByID(bucketName, id string, record interface{}) error {
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		recordJSON := b.Get([]byte(id))
		if recordJSON == nil {
			return fmt.Errorf("item with ID %s not found", id) // Return error if item is not found
		}

		if err := json.Unmarshal(recordJSON, record); err != nil {
			return err
		}

		// Decrypt the content if the record type is models.Item
		if item, ok := record.(*models.Item); ok {
			decodedBytes,
				err := base64.StdEncoding.DecodeString(
				item.Content.Description,
			)

			if err != nil {
				return err
			}

			decryptedContent,
				err := cryptography.DecryptData(
				decodedBytes,
				config.Env(exports.EnvPath),
			)

			if err != nil {
				return err
			}

			item.Content.Description = string(decryptedContent)

			decodedBytes,
				err = base64.StdEncoding.DecodeString(
				item.Content.Image,
			)

			if err != nil {
				return err
			}

			decryptedContent, err = cryptography.DecryptData(
				decodedBytes,
				config.Env(exports.EnvPath),
			)

			if err != nil {
				return err
			}
			item.Content.Image = string(decryptedContent)
		}

		return nil
	})
}

// GetUserByUsername retrieves a user by username from the database
func GetUserByUsername(username string) (*models.User, error) {
	return GetUserByField("user", username)
}

// GetUserByWallet retrieves a user by wallet address from the database
func GetUserByWallet(wallet string) (*models.User, error) {
	return GetUserByField("wallet", wallet)
}

// getUserByField retrieves a user by a specific field (e.g., username or wallet) from the database
func GetUserByField(field string, value string) (*models.User, error) {
	var user models.User
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return fmt.Errorf("bucket %q not found ", "users")
		}

		// Iterate through the bucket to find the user by the specified field
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var u models.User
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			if field == "user" && u.User == value {
				user = u
				return nil
			} else if field == "wallet" && u.Wallet == value {
				user = u
				return nil
			}
		}
		return nil // User not found
	})

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil // User not found
	}

	return &user, nil
}

func DeleteRecord(bucketName, id string) error {
	return db.Update(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %q not found ", bucketName)
			}

			if err := b.Delete([]byte(id)); err != nil {
				return fmt.Errorf("failed to delete item with ID %s: %w", id, err)
			}

			return nil
		},
	)
}

func NextID(bucketName string) (int, error) {
	var id int
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %q not found ", bucketName)
		}

		// Get the current sequence number
		seq, err := b.NextSequence()
		if err != nil {
			return err
		}

		id = int(seq)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
