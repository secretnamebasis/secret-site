package models

import "time"

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
