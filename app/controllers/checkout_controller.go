package controllers

import (
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

func CreateUserCheckout(order *models.JSON_User_Order) (checkout models.Checkout, err error) {
	if err = order.Validate(); err != nil {
		return checkout, err
	}
	var timestamp time.Time = time.Now()
	var expiration time.Time = time.Now().Add(5 * time.Minute)
	destination, err := dero.GetWalletAddress(config.WalletEndpoint)
	if err != nil {
		return checkout, err
	}
	a := rpc.Arguments{
		rpc.Argument{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    order.Name,
		},
		rpc.Argument{
			Name:     rpc.RPC_VALUE_TRANSFER,
			DataType: rpc.DataUint64,
			Value:    config.UserRegistrationFee,
		},
		// rpc.Argument{ // this doesn't work as expected
		// 	Name:     rpc.RPC_EXPIRY,
		// 	DataType: rpc.DataTime,
		// 	Value:    expiration,
		// },
		rpc.Argument{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    config.UserRegistrationPort,
		},
		rpc.Argument{
			Name:     rpc.RPC_NEEDS_REPLYBACK_ADDRESS,
			DataType: rpc.DataUint64,
			Value:    uint64(1),
		},
	}
	p := rpc.Make_Integrated_Address_Params{
		Address:     destination.String(),
		Payload_RPC: a,
	}
	addr, err := dero.MakeIntegratedAddress(p)
	if err != nil {
		return checkout, err
	}

	id, err := NextCheckoutID()
	if err != nil {
		return checkout, err
	}
	checkout = models.Checkout{
		ID:         id,
		Address:    addr.Integrated_Address,
		CreatedAt:  timestamp,
		Expiration: expiration,
	}

	if err := checkout.Validate(); err != nil {
		return checkout, err
	}

	checkout.Initialize()

	err = database.CreateRecord(bucketCheckouts, &checkout)
	return checkout, err
}

// NextUserID returns the next available user ID.
func NextCheckoutID() (int, error) {
	return database.NextID(bucketCheckouts)
}
