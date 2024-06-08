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
	if i.Title == "" || i.Description == "" || i.SCID == "" {
		return errors.New("cannot be empty item json")
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
func (u *JSON_User_Order) Validate() error {
	if u.Name == "" {
		return errors.New("name and wallet cannot be empty user json")
	}
	if len(u.Name) > // also you are being lazy, we all know that isn't 66 characters
		66 { // this is an arbitrary number (tbh, it is probably still too big)
		return errors.New("username is longer than 66 characters")
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
