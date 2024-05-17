package models

import (
	"errors"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Wallet string `json:"wallet"`
	// Password   []byte    `json:"password"`
	Role       string    `json:"role"`
	LastSignIn time.Time `json:"last_sign_in"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NewUser creates a new User instance with the provided data
func (u *User) Initialize() *User {
	// Generate ID and password

	return &User{
		ID:     u.ID,
		Name:   u.Name,
		Wallet: u.Wallet,
		// Password:  u.Password,
		Role:      "user",
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Validate method validates the user data.
func (u *User) Validate() error {

	if err := u.isEmpty(); err != nil {
		return errors.New("submission is empty")
	}
	if err := hasValidWallet(u.Wallet); err != nil {
		return errors.New("invalid wallet address")
	}

	// we are no longer going to use passwords
	// if u.Password == nil {
	// 	return errors.New("password cannot be empty")
	// }
	// Add more validation rules as needed
	return nil
}

// hasValidWallet checks if the provided wallet address is valid
func hasValidWallet(wallet string) error {
	// Attempts to fetch the encrypted balance of the wallet address
	_, err := dero.GetEncryptedBalance(config.NodeEndpoint, wallet)
	if err != nil {
		return err
	}
	return nil
}

// validateUserData checks if the provided user data is valid
func (u *User) isEmpty() error {
	if u.Name == "" ||
		u.Wallet == "" ||
		// we aren't using passwords anymore
		// the validation is from the async operation from the wallet
		// u.Password == nil ||
		u.ID == 0 {
		return errors.New("user and wallet fields are required")
	}
	return nil
}
