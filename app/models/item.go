package models

import (
	"errors"
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
)

// Item represents a sample data structure for demonstration
type Item struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	SCID      string    `json:"scid"`
	Data      []byte    `json:"data"` // ItemData
	ImageURL  string    `json:"image_url"`
	FileURL   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemData struct {
	Description string `json:"description"`
	Image       string `json:"image"`
	File        string `json:"file"`
}

// InitializeItem creates and initializes a new Item instance
func (i *Item) Initialize() *Item {
	timestamp := time.Now()
	i.ImageURL = config.Domainname + "/images/" + i.SCID
	i.FileURL = config.Domainname + "/files/" + i.SCID
	i.CreatedAt = timestamp
	i.UpdatedAt = timestamp
	item := &Item{
		ID:        i.ID,
		Title:     i.Title,
		SCID:      i.SCID,
		Data:      []byte{},
		ImageURL:  i.ImageURL,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}

	return item
}

// Validate method validates the fields of the Item struct
func (i *Item) Validate() error {
	if i.Data == nil ||
		i.ID == 0 ||
		i.CreatedAt == (time.Time{}) ||
		i.UpdatedAt == (time.Time{}) {

		return errors.New("cannot be empty")
	}

	return nil
}
