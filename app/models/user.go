package models

import "time"

type User struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Wallet    string    `json:"wallet"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
