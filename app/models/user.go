package models

import (
	"errors"
	"time"

	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

type User struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Wallet    string    `json:"wallet"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate method validates the user data.
func (u *User) Validate() error {
	if err := isEmpty(u); err != nil {

		return errors.New("invalid wallet address")

	}
	if err := isValidWallet(u.Wallet); err != nil {
		return errors.New("invalid wallet address")
	}
	return nil
	// when creating a user... they don't have a password on file yet
	// if u.Password == "" {
	// 	return errors.New("password cannot be empty")
	// }
	// Add more validation rules as needed

}

// isValidWallet checks if the provided wallet address is valid
func isValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(wallet)
	return err
}

// validateUserData checks if the provided user data is valid
func isEmpty(user *User) error {
	if user.User == "" || user.Wallet == "" {
		return errors.New("user and wallet fields are required")
	}
	return nil
}
