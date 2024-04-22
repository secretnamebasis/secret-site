package models

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

type User struct {
	// ID represents the unique identifier of the user.
	ID int `json:"id"`
	// Name stores the name of the user.
	Name string `json:"name"`
	// Wallet stores the DERO wallet address of the user.
	Wallet string `json:"wallet"`
	// Password stores the hashed password of the user.
	Password []byte `json:"password"`
	// Role represents the roles assigned to the user.
	Role []string `json:"roles"`
	// LastSignIn stores the timestamp of the user's last sign-in.
	LastSignIn time.Time `json:"last_sign_in"`
	// CreatedAt stores the timestamp when the user was created.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt stores the timestamp when the user was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new User instance with the provided data
func (u *User) Initialize() *User {
	// Generate ID and password

	return &User{
		ID:        u.ID,
		Name:      u.Name,
		Wallet:    u.Wallet,
		Password:  u.Password,
		Role:      []string{"user"},
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

	// when creating a user... they don't have a password on file yet
	if u.Password == nil {
		return errors.New("password cannot be empty")
	}
	// Add more validation rules as needed
	return nil
}

// hasValidWallet checks if the provided wallet address is valid
func hasValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(config.NodeEndpoint, wallet)
	if err != nil {
		log.Errorf("reg: %s", err)
	}
	return err
}

// validateUserData checks if the provided user data is valid
func (u *User) isEmpty() error {
	if u.Name == "" ||
		u.Wallet == "" ||
		u.ID == 0 ||
		u.Password == nil {
		return errors.New("user and wallet fields are required")
	}
	return nil
}
