package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
)

// Item represents a sample data structure for demonstration
type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	// ItemData
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageURL  string
}
type ItemData struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

// InitializeItem creates and initializes a new Item instance
func (i *Item) Initialize() *Item {

	item := &Item{
		ID:        i.ID,
		Title:     i.Title,
		Data:      []byte{},
		ImageURL:  config.Domainname + "/images/" + fmt.Sprintf("%d", i.ID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return item
}

// Validate method validates the fields of the Item struct
func (i *Item) Validate() error {
	if i.Data == nil || i.ID == 0 {
		return errors.New("cannot be empty")
	}

	return nil
}
