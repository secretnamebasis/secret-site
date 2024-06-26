package models

import (
	"errors"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
)

type JSON_Item_Order struct {
	Title       string          `json:"title"`
	SCID        string          `json:"scid"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	File        string          `json:"file"`
	User        JSON_User_Order `json:"user"`
}

// Validate method validates the fields of the Item struct
func (i *JSON_Item_Order) Validate() error {
	if i.Title == "" || i.Description == "" || i.SCID == "" || i.User == (JSON_User_Order{}) {
		return errors.New("cannot be empty")
	}
	if err := hasValidSCID(i.SCID); err != nil {
		return err
	}
	return nil
}

type JSON_User_Order struct {
	Name     string `json:"name"`
	Wallet   string `json:"wallet"`
	Password string `json:"password"`
}

// Validate method validates the fields of the JSON_User_Order struct
func (i *JSON_User_Order) Validate() error {
	if i.Name == "" || i.Wallet == "" {
		return errors.New("name and wallet cannot be empty")
	}
	if err := hasValidWallet(i.Wallet); err != nil {
		return err
	}
	return nil
}

// hasValidSCID checks if the provided SCID is valid
func hasValidSCID(scid string) error {
	// Attempt to fetch the code of the contract
	result, err := dero.GetSCID(config.NodeEndpoint, scid)
	if err != nil {
		return err
	}

	if result.Code == "" {
		return errors.New("error: code is empty")
	}
	return nil
}
