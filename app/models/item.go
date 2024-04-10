package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/secretnamebasis/secret-site/app/exports"
)

// Item represents a sample data structure for demonstration
type Item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   Content   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Content struct {
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageURL    string
}

// InitializeItem creates and initializes a new Item instance
func InitializeItem(
	id int,
	title, description, image string,
) *Item {
	item := &Item{
		ID:    id,
		Title: title,
		Content: Content{
			Description: description,
			Image:       image,
			ImageURL:    exports.DOMAINNAME + "/images/" + fmt.Sprintf("%d", id),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return item
}

// Validate method validates the fields of the Item struct
func (i *Item) Validate() error {
	if i.Title == "" {
		return errors.New("title cannot be empty")
	}
	if i.Content.Description == "" {
		return errors.New("content cannot be empty")
	}
	// Add more validation rules as needed
	// Validate Content

	return nil
}
