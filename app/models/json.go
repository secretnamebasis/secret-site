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

// Validate method validates the fields of the Item struct
func (i *JSON_User_Order) Validate() error {
	if i.Name == "" || i.Wallet == "" {
		return errors.New("cannot be empty")
	}
	return nil
}
