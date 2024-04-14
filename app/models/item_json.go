package models

import "errors"

type JSON_Item_Order struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// Validate method validates the fields of the Item struct
func (i *JSON_Item_Order) Validate() error {
	if i.Title == "" || i.Description == "" {
		return errors.New("cannot be empty")
	}
	return nil
}
