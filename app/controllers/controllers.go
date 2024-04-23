package controllers

import (
	"errors"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

// Define bucket names
const (
	bucketItems = "items"
	bucketUsers = "users"
)

// isValidWallet checks if the provided wallet address is valid
func isValidWallet(wallet string) error {
	// Attempt to fetch the balance of the wallet address
	_, err := dero.GetEncryptedBalance(config.NodeEndpoint, wallet)
	return err
}

// validateWalletAddress checks if the provided wallet address is valid
func ValidateWalletAddress(wallet string) error {
	if err := isValidWallet(wallet); err != nil {
		return errors.New("invalid wallet address")
	}
	return nil
}
