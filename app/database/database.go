package database

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/cryptography"
	"github.com/secretnamebasis/secret-site/app/models"

	"go.etcd.io/bbolt"
)

var (
	db          *bbolt.DB
	itemsBucket = []byte("items")
	usersBucket = []byte("users")
)

func Initialize(c config.Server) error {
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
			// just so it doesn't get lost in the shuffle
			var id int

			switch r := record.(type) {
			case *models.Item: // Ensure we're dealing with a pointer to models.Item
				id = r.ID
				// Encrypt content before storing in the database
				encryptedBytes, err := cryptography.EncryptData(
					r.Data, // we want to lock these bitches down!
					config.Env( // and to do it we are going into our env
						"SECRET", // and we are going to refer to our secret
					), // shouldn't ever change unless we are changing our encryption scheme
				)
				// and it better work.
				if err != nil {
					return err
				}
				r.Data = encryptedBytes
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
						// Create a new instance of the record type
						newRecord := reflect.New(reflect.TypeOf(recordType).Elem()).Interface()

						// Unmarshal the JSON data into the new record
						if err := json.Unmarshal(v, newRecord); err != nil {
							return err
						}

						// Decrypt the content if the record type is models.Item
						if item, ok := newRecord.(*models.Item); ok {
							// Decrypt the content before assigning it to the record
							_, err := cryptography.DecryptData(
								item.Data,
								config.Env("SECRET"),
							)
							if err != nil {
								return err
							}

						}

						// Append the record to the slice
						sliceValue := reflect.ValueOf(records).Elem()
						sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(newRecord).Elem()))

						return nil
					},
				)
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

		switch updatedRecord := updatedRecord.(type) {
		case *models.JSONItemData:
			var existingItem models.Item
			if err := json.Unmarshal(recordJSON, &existingItem); err != nil {
				return err
			}

			decryptedData, err := cryptography.DecryptData(existingItem.Data, config.Env("SECRET"))
			if err != nil {
				return err
			}

			var existingItemData models.ItemData
			if err := json.Unmarshal(decryptedData, &existingItemData); err != nil {
				return err
			}
			if updatedRecord.Title != "" {
				existingItem.Title = updatedRecord.Title
			}
			// Update the existingItemData fields
			if updatedRecord.Image != "" {
				existingItemData.Image = updatedRecord.Image
			}
			if updatedRecord.Description != "" {
				existingItemData.Description = updatedRecord.Description
			}

			// Marshal the updated data and encrypt it
			updatedBytes, err := json.Marshal(existingItemData)
			if err != nil {
				return err
			}

			encryptedBytes, err := cryptography.EncryptData(updatedBytes, config.Env("SECRET"))
			if err != nil {
				return err
			}

			// Update existingItem with the encrypted data and set the updated timestamp
			existingItem.Data = encryptedBytes
			existingItem.UpdatedAt = time.Now()

			// Marshal the updated existingItem and store it in the bucket
			updatedItemJSON, err := json.Marshal(existingItem)
			if err != nil {
				return err
			}
			return b.Put([]byte(id), updatedItemJSON)

		case *models.User:
			var existingUser models.User
			if err := json.Unmarshal(recordJSON, &existingUser); err != nil {
				return err
			}
			updateExistingUser(updatedRecord, &existingUser)
			// Marshal the updated existingUser and store it in the bucket
			updatedUserJSON, err := json.Marshal(existingUser)
			if err != nil {
				return err
			}
			return b.Put([]byte(id), updatedUserJSON)

		default:
			return fmt.Errorf("unsupported record type")
		}
	})
}

func checkError(err error) error {
	if err != nil {
		return err
	}
	return nil
}

// updateRecordFromJSON updates the provided record using data from JSON.
func updateRecordFromJSON(recordJSON []byte, updatedRecord interface{}) error {
	switch record := updatedRecord.(type) {
	case *models.JSONItemData:
		var existingItem models.Item
		if err := json.Unmarshal(recordJSON, &existingItem); err != nil {
			return err
		}
		fmt.Printf("EXISTING BYTES: %s\n", existingItem.Data)

		decryptedData, err := cryptography.DecryptData(existingItem.Data, config.Env("SECRET"))
		fmt.Printf("DECRYPTED BYTES: %s\n", decryptedData)
		if err != nil {
			return err
		}

		var existingItemData models.ItemData

		// Fill the existingItemData with the decryptedData
		if err := json.Unmarshal(decryptedData, &existingItemData); err != nil {
			return err
		}
		fmt.Printf("Unmarshalled Data: %+v\n", existingItemData)
		fmt.Printf("Record: %s\n", record)

		// Define a map to store the fields to be updated
		fieldsToUpdate := map[string]bool{
			"Data.Image":       record.Image != "",
			"Data.Description": record.Description != "",
		}

		// Update the fields in the existingItemData
		if fieldsToUpdate["Data.Image"] {
			existingItemData.Image = record.Image
		}
		if fieldsToUpdate["Data.Description"] {
			existingItemData.Description = record.Description
		}

		updatedBytes, err := json.Marshal(existingItemData)
		fmt.Printf("UPDATED BYTES: %s\n", updatedBytes)
		if err != nil {
			return err
		}

		encryptedBytes, err := cryptography.EncryptData(updatedBytes, config.Env("SECRET"))
		fmt.Printf("Encrypted BYTES: %s\n", encryptedBytes)
		if err != nil {
			return err
		}

		existingItem.Data = encryptedBytes
		existingItem.UpdatedAt = time.Now()

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

// updateUserFromExisting updates the user record based on existing data.
func updateExistingUser(updatedUser, existingUser *models.User) {
	if updatedUser.Name != "" {
		existingUser.Name = updatedUser.Name
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
	return db.View(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %s not found", bucketName)
			}

			recordJSON := b.Get([]byte(id))
			if recordJSON == nil {
				return fmt.Errorf("record with ID %s not found", id) // Return error if item is not found
			}

			if err := json.Unmarshal(recordJSON, record); err != nil {
				return err
			}

			// Decrypt the content if the record type is models.Item
			if item, ok := record.(*models.Item); ok {
				fmt.Printf("DATA: %s\n", item.Data)
				// Marshal item's Data field to JSON

				// Encrypt JSON bytes
				decryptedBytes, err := cryptography.DecryptData(item.Data, config.Env("SECRET"))
				if err != nil {
					return err
				}
				fmt.Printf("DECRYPTED DATA: %s\n", decryptedBytes)

				item.Data = decryptedBytes
			}

			return nil
		},
	)
}

func GetItemByTitle(title string) (*models.Item, error) {
	return GetItemByField("title", title)
}

// getUserByField retrieves a user by a specific field (e.g., username or wallet) from the database
func GetItemByField(field string, value string) (*models.Item, error) {
	var item models.Item
	err := db.View(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("items"))
			if b == nil {
				return fmt.Errorf("bucket %q not found ", "items")
			}

			// Iterate through the bucket to find the user by the specified field
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var i models.Item
				// and then you are going to take that user object
				// and you are going to let it act as sorting
				// bins for the bytes of data that you are
				// going to spill out all over the place
				if err := json.Unmarshal(v, &i); err != nil {
					return err
				}
				if field == "title" && i.Title == value {
					item = i
					return nil
				} else if field == "data" && string(i.Data) == value {
					item = i
					return nil
				}
			}
			return nil // User not found
		},
	)

	if err != nil {
		return nil, err
	}

	if item.ID == 0 {
		return nil, nil // item not found
	}
	// now return the object back upstairs.
	return &item, nil
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
	err := db.View(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("users"))
			if b == nil {
				return fmt.Errorf("bucket %q not found ", "users")
			}

			// Iterate through the bucket to find the user by the specified field
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var u models.User
				// and then you are going to take that user object
				// and you are going to let it act as sorting
				// bins for the bytes of data that you are
				// going to spill out all over the place
				if err := json.Unmarshal(v, &u); err != nil {
					return err
				}
				if field == "user" && u.Name == value {
					user = u
					return nil
				} else if field == "wallet" && u.Wallet == value {
					user = u
					return nil
				}
			}
			return nil // User not found
		},
	)

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, nil // User not found
	}
	// now return the object back upstairs.
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
	err := db.Update(
		func(tx *bbolt.Tx) error {
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
