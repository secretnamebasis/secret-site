package models

import "errors"

type JSON_Item_Order struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	User        JSON_User_Order `json:"user"`
}

// Validate method validates the fields of the Item struct
func (i *JSON_Item_Order) Validate() error {
	if i.Title == "" || i.Description == "" || i.User == (JSON_User_Order{}) {
		return errors.New("cannot be empty")
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
