package db

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/models"

	"go.etcd.io/bbolt"
)

var (
	DB          *bbolt.DB
	itemsBucket = []byte("items")
	usersBucket = []byte("users")
)

func InitDB(env string) error {
	var databasePath = "database/" + env + ".db"
	var err error
	if err := os.MkdirAll("database", 0755); err != nil {
		return err
	}

	DB, err = bbolt.Open(databasePath, 0600, nil)
	if err != nil {
		return err
	}

	// Create buckets if they don't exist
	err = DB.Update(func(tx *bbolt.Tx) error {
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
	return DB.Update(
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
				encrypted, err := cryptography.EncryptData([]byte(r.Content), config.Env("SECRET"))
				if err != nil {
					return err
				}

				// Encode the encrypted content to Base64
				encryptedContent = base64.StdEncoding.EncodeToString(encrypted)
				r.Content = encryptedContent

			case *models.User:
				id = r.ID

			default:
				return fmt.Errorf("unsupported record type")
			}

			recordJSON, err := json.Marshal(record)
			if err != nil {
				return err
			}
			log.Printf("Record JSON: %s", recordJSON) // Log record JSON
			return b.Put([]byte(strconv.Itoa(id)), recordJSON)
		},
	)
}

// GetAllRecords retrieves all records from the specified bucket and unmarshals them into the provided slice.
func GetAllRecords(bucketName string, records interface{}, c *fiber.Ctx) error {
	return DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %q not found", bucketName)
		}

		// Define a helper function to unmarshal records
		unmarshalRecord := func(recordType interface{}) error {
			return b.ForEach(func(k, v []byte) error {
				if err := json.Unmarshal(v, recordType); err != nil {
					return err
				}
				// Decrypt the content if the record type is models.Item
				if item, ok := recordType.(*models.Item); ok {
					decodedBytes, err := base64.StdEncoding.DecodeString(item.Content)
					if err != nil {
						return err
					}
					decryptedContent, err := cryptography.DecryptData(decodedBytes, config.Env("SECRET"))
					if err != nil {
						return err
					}
					item.Content = string(decryptedContent)
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

// GetRecordByID retrieves a record from the specified bucket by ID and unmarshals it into the provided model.
func GetRecordByID(bucketName, id string, record interface{}) error {
	return DB.View(func(tx *bbolt.Tx) error {
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
			decodedBytes, err := base64.StdEncoding.DecodeString(item.Content)
			if err != nil {
				return err
			}
			decryptedContent, err := cryptography.DecryptData(decodedBytes, config.Env("SECRET"))
			if err != nil {
				return err
			}
			item.Content = string(decryptedContent)
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
	err := DB.View(func(tx *bbolt.Tx) error {
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

// UpdateRecord updates a record in the specified bucket with the provided ID and updated data.
func UpdateRecord(bucketName, id string, updatedRecord interface{}) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucketName)
		}

		recordJSON := b.Get([]byte(id))
		if recordJSON == nil {
			return fmt.Errorf("record with ID %s not found in bucket %s", id, bucketName)
		}
		switch updatedRecord := updatedRecord.(type) {
		case *models.Item:
			var existingItem models.Item
			if err := json.Unmarshal(recordJSON, &existingItem); err != nil {
				return err
			}
			// Update only the non-zero fields
			if updatedRecord.Title != "" {
				existingItem.Title = updatedRecord.Title
			}
			if updatedRecord.Content != "" {
				existingItem.Content = updatedRecord.Content
			}
			// Preserve the ID and creation timestamp
			updatedRecord.ID = existingItem.ID
			updatedRecord.CreatedAt = existingItem.CreatedAt
			updatedRecord.UpdatedAt = time.Now()
		case *models.User:
			var existingUser models.User
			if err := json.Unmarshal(recordJSON, &existingUser); err != nil {
				return err
			}
			// Update only the non-zero fields
			if updatedRecord.User != "" {
				existingUser.User = updatedRecord.User
			}
			if updatedRecord.Wallet != "" {
				existingUser.Wallet = updatedRecord.Wallet
			}
			// Always update the password if provided
			if updatedRecord.Password != "" {
				existingUser.Password = updatedRecord.Password
			} else {
				// If password is not provided, preserve the existing password
				updatedRecord.Password = existingUser.Password
			}
			// Preserve the ID and creation timestamp
			updatedRecord.ID = existingUser.ID
			updatedRecord.CreatedAt = existingUser.CreatedAt
			updatedRecord.UpdatedAt = time.Now()
		default:
			return fmt.Errorf("unsupported record type")
		}

		updatedJSON, err := json.Marshal(updatedRecord)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), updatedJSON)
	})
}

func DeleteRecord(bucketName, id string) error {
	return DB.Update(
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
	err := DB.Update(func(tx *bbolt.Tx) error {
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
