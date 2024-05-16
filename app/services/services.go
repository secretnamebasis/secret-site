package services

import (
	"fmt"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

func ProcessCheckouts(c config.Server) error {
	for {
		// so let's go get our checkouts
		var checkouts []models.Checkout
		var err = database.GetAllRecords("checkouts", &checkouts)
		if err != nil {
			return err
		}
		// and lets range through each one
		for _, checkout := range checkouts {

			fmt.Println(checkout)

			// older than 5 min checkouts should be deleted
			if checkout.Expiration.Unix() < time.Now().Add(-5*time.Minute).Unix() {
				strID := fmt.Sprint(checkout.ID)
				err = database.DeleteRecord(
					"checkouts",
					strID,
				)
				if err != nil {
					return err
				}
				continue
			} else // let's go ahead and process each checkout
			{
				// so lets establish that this checkout is an address
				addr, err := rpc.NewAddress(checkout.Address)
				if err != nil {
					return err
				}

				if !addr.Arguments.Has(
					rpc.RPC_DESTINATION_PORT,
					rpc.DataUint64,
				) ||
					!addr.Arguments.Has(
						rpc.RPC_COMMENT,
						rpc.DataString,
					) {
					continue
				}

				// then define the port used
				addrPort := addr.Arguments.Value(
					rpc.RPC_DESTINATION_PORT,
					rpc.DataUint64,
				).(uint64)

				// and the comment encoded
				addrComment := addr.Arguments.Value(
					rpc.RPC_COMMENT,
					rpc.DataString,
				).(string)

				// let's set our params
				params := rpc.Get_Transfers_Params{
					DestinationPort: addrPort,        // at our desired port
					In:              true,            // only incoming transfers
					SCID:            crypto.ZEROHASH, // DERO, but could be any scid...
					Out:             false,
					Coinbase:        false,
				}
				// then get the transfers fromt the wallet
				transfers, err := dero.GetWalletTransfers(
					c.WalletEndpoint,
					params,
				)

				if err != nil {
					return err
				}

				// now lets range through the entries
				for _, transfer := range transfers.Entries {
					if !transfer.Payload_RPC.Has(
						rpc.RPC_COMMENT,
						rpc.DataString,
					) {
						continue
					}

					// define the transfer's comment
					txComment := transfer.Payload_RPC.Value(
						rpc.RPC_COMMENT,
						rpc.DataString,
					).(string)

					// as well as its destination port
					txPort := transfer.DestinationPort

					if txPort != addrPort ||
						txComment != addrComment {
						continue
					}
					fmt.Println(transfer)
					fmt.Println("DO SOMETHING HERE!!!")

				}
			}
		}
		time.Sleep(5 * time.Second)

	}
	return nil
}
