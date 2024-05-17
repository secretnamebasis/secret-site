package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/models"

	"go.etcd.io/bbolt"
)

var (
	db                                       *bbolt.DB
	itemsBucket, usersBucket, checkoutBucket = []byte("items"), []byte("users"), []byte("checkouts")

	// this was my first byte array.
	buckets = [][]byte{
		itemsBucket,
		checkoutBucket,
		usersBucket,
	}
)

func Initialize(c config.Server) error {
	// Set directory of the database
	if err := os.MkdirAll(c.DatabasePath, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(
		c.DatabasePath,
		c.Environment+".db",
	)

	// Open or create the database file
	var err error
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}

	// Ensure buckets exist
	err = db.Update(
		func(tx *bbolt.Tx) error {
			for _, bucket := range buckets {
				_, err := tx.CreateBucketIfNotExists(bucket)
				if err != nil {
					return err
				}
			}
			return nil
		})

	return err
}

// CreateRecord creates a record in the database for the given item/user after encrypting the content.
func CreateRecord(bucketName string, record any) error {
	return db.Update(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %q not found ", bucketName)
			}

			value := reflect.ValueOf(record)
			id := value.Elem().FieldByName("ID")
			if !id.IsValid() {
				return fmt.Errorf("record does not have ID field")
			}

			i := int(id.Int())

			recordJSON, err := json.Marshal(record)
			if err != nil {
				return err
			}

			return b.Put([]byte(strconv.Itoa(i)), recordJSON)
		},
	)
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

			return nil
		},
	)
}

// GetAllRecords retrieves all records from the specified bucket and unmarshals them into the provided slice.
func GetAllRecords(bucketName string, records any) error {
	return db.View(
		func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return fmt.Errorf("bucket %q not found", bucketName)
			}

			// Define a helper function to unmarshal records
			unmarshalRecord := func(recordType any) error {
				return b.ForEach(
					func(k, v []byte) error {
						r := reflect.TypeOf(recordType).Elem()
						// Create a new instance of the record type
						newRecord := reflect.New(r).Interface()

						// Unmarshal the JSON data into the new record
						if err := json.Unmarshal(v, newRecord); err != nil {
							return err
						}

						// Append the record to the slice
						sliceValue := reflect.ValueOf(records).Elem()
						sliceValue.Set(
							reflect.Append(
								sliceValue,
								reflect.ValueOf(
									newRecord,
								).Elem(), // which is an interface
							),
						)

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
			case *[]models.Checkout:
				return unmarshalRecord(&models.Checkout{})
			default:
				return fmt.Errorf("unsupported record type")
			}
		})
}

// GetAllRecords retrieves all records from the specified bucket and unmarshals them into the provided slice.
func GetAllItemTitles(bucketName string, records interface{}) error {
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
						reflectRecord := reflect.TypeOf(recordType).Elem()
						// Create a new instance of the record type
						newRecord := reflect.New(
							reflectRecord,
						).Interface()

						// Unmarshal the JSON data into the new record
						if err := json.Unmarshal(v, newRecord); err != nil {
							return err
						}

						// Append the record to the slice
						sliceValue := reflect.ValueOf(records).Elem()
						reflectNewRecord := reflect.ValueOf(newRecord).Elem()
						sliceValue.Set(reflect.Append(
							sliceValue,
							reflectNewRecord,
						),
						)

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

func GetItemByTitle(title string) (models.Item, error) {
	return GetItemByField("title", title)
}

// getUserByField retrieves a user by a specific field (e.g., username or wallet) from the database
func GetItemByField(field string, value string) (models.Item, error) {
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
				} else if field == "scid" && i.SCID == value {
					item = i
					return nil
				}
			}
			return nil // User not found
		},
	)

	if err != nil {
		return item, err
	}

	if item.ID == 0 {
		return item, nil // item not found
	}
	// now return the object back upstairs.
	return item, nil
}

// GetUserByUsername retrieves a user by username from the database
func GetUserByUsername(username string) (models.User, error) {
	return GetUserByField("user", username)
}

// GetUserByWallet retrieves a user by wallet address from the database
func GetUserByWallet(wallet string) (models.User, error) {
	return GetUserByField("wallet", wallet)
}

// getUserByField retrieves a user by a specific field (e.g., username or wallet) from the database
func GetUserByField(field string, value string) (models.User, error) {
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
		return user, err
	}
	// here we are running into a problem.
	// we don't have an error, but we don't have a valid user either
	if user.ID == 0 {
		return user, nil // User not found
	}
	// now return the object back upstairs.
	return user, nil
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
			return fmt.Errorf("bucket %q not found", bucketName)
		}

		// Check if the bucket is empty
		stats := b.Stats()
		if stats.KeyN == 0 {
			// If the bucket is empty, set the sequence number to 1
			if err := b.SetSequence(1); err != nil {
				return err
			}
			id = 1
			return nil
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
