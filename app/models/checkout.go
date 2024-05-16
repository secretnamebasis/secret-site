package models

import (
	"errors"
	"time"
)

type Checkout struct {
	ID         int
	Address    string
	CreatedAt  time.Time
	Expiration time.Time
}

func (c *Checkout) Initialize() *Checkout {
	return &Checkout{
		ID:         c.ID,
		Address:    c.Address,
		CreatedAt:  c.CreatedAt,
		Expiration: c.Expiration,
	}
}

func (c *Checkout) Validate() error {
	if c.ID == 0 ||
		c.Address == "" ||
		c.CreatedAt == (time.Time{}) ||
		c.Expiration == (time.Time{}) {

		return errors.New("cannot be empty")
	}

	return nil
}
