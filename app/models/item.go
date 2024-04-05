package models

import (
	"errors"
	"time"
)

// Item represents a sample data structure for demonstration
type Item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate method validates the fields of the Item struct
func (i *Item) Validate() error {
	if i.Title == "" {
		return errors.New("title cannot be empty")
	}
	if i.Content == "" {
		return errors.New("content cannot be empty")
	}
	// Add more validation rules as needed

	return nil
}
