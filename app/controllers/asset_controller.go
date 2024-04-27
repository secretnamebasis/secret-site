package controllers

import (
	"time"

	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateAssetCheckout(
	order *models.JSON_Asset_Order,
) (
	*models.Checkout,
	error,
) {
	var expiry time.Time = time.Now().Add(5 * time.Minute)
	var price uint64 = 10000 // honestly this should be dynamicly determined

	address, err := dero.MakeIntegratedAddress(
		order.Name,
		price,
		expiry,
	)
	if err != nil {
		return nil, err
	}
	var checkout models.Checkout
	id, err := NextCheckoutID()
	if err != nil {
		return &models.Checkout{}, err
	}
	checkout.ID = id
	checkout.Address = address.Integrated_Address
	checkout.Expiration = expiry
	checkout.Initialize()
	if err := database.CreateRecord(
		bucketCheckouts,
		&checkout,
	); err != nil {
		return &models.Checkout{}, err
	}
	return &checkout, err
}

func CreateAsset(order models.JSON_Asset_Order) {

	// Create the item record
	_, _ = dero.MintContract(
		config.WalletEndpoint,
		dero.NFAContract(
			order.Royalty,
			order.Name,
			order.Description,
			order.Type,
			order.Collection,
			order.Wallet,
		),
		order.Wallet,
	)
}

// NextCheckoutID returns the next available checkout ID.
func NextCheckoutID() (int, error) {
	return database.NextID(bucketCheckouts)
}
