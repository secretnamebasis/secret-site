package controllers

import (
	"fmt"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateItemCheckout(order *models.JSON_Item_Order) (models.Checkout, error) {
	return createCheckout(
		order,
		order.Title,
		config.ItemListingFee,
		config.ItemListingPort,
	)
}

func CreateUserCheckout(order *models.JSON_User_Order) (models.Checkout, error) {
	return createCheckout(
		order,
		order.Name,
		config.UserRegistrationFee,
		config.UserRegistrationPort,
	)
}
func createCheckout(
	order interface{},
	comment string,
	valueTransfer uint64,
	destinationPort uint64,
) (models.Checkout, error) {
	var checkout models.Checkout
	var errors []error

	// Validate order
	switch o := order.(type) {
	case *models.JSON_Item_Order:
		if err := o.Validate(); err != nil {
			errors = append(errors, err)
		}
	case *models.JSON_User_Order:
		if err := o.Validate(); err != nil {
			errors = append(errors, err)
		}
	default:
		return checkout, fmt.Errorf("unsupported order type")
	}

	timestamp := time.Now()
	expiration := timestamp.Add(5 * time.Minute)

	// Get wallet address
	destination, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		errors = append(errors, err)
	}

	// Prepare arguments for integrated address
	a := rpc.Arguments{
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    comment,
		},
		{
			Name:     rpc.RPC_VALUE_TRANSFER,
			DataType: rpc.DataUint64,
			Value:    valueTransfer,
		},
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    destinationPort,
		},
		{
			Name:     rpc.RPC_NEEDS_REPLYBACK_ADDRESS,
			DataType: rpc.DataUint64,
			Value:    uint64(1),
		},
	}

	p := rpc.Make_Integrated_Address_Params{
		Address:     destination.String(),
		Payload_RPC: a,
	}

	// Make integrated address
	addr, err := dero.MakeIntegratedAddress(p)
	if err != nil {
		errors = append(errors, err)
	}

	// Generate next checkout ID
	id, err := NextCheckoutID()
	if err != nil {
		errors = append(errors, err)
	}

	// Create checkout object
	checkout = models.Checkout{
		ID:         id,
		Address:    addr.Integrated_Address,
		CreatedAt:  timestamp,
		Expiration: expiration,
	}

	// Validate checkout
	if err := checkout.Validate(); err != nil {
		errors = append(errors, err)
	}

	// Initialize checkout
	checkout.Initialize()

	// Create record in database
	if err := database.CreateRecord(bucketCheckouts, &checkout); err != nil {
		errors = append(errors, err)
	}

	// Check for errors
	if len(errors) > 0 {
		return checkout, errors[0]
	}

	return checkout, nil
}

// NextUserID returns the next available user ID.
func NextCheckoutID() (int, error) {
	return database.NextID(bucketCheckouts)
}
