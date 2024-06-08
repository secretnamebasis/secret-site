package services

import (
	"fmt"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/controllers"
	"github.com/secretnamebasis/secret-site/app/database"
	"github.com/secretnamebasis/secret-site/app/integrations/dero"
	"github.com/secretnamebasis/secret-site/app/models"
)

func ProcessCheckouts(c config.Server) error {
	go func() error {

		for {
			// so let's go get our checkouts
			var checkouts []models.Checkout
			var err = database.GetAllRecords("checkouts", &checkouts)
			if err != nil {
				return err
			}
			// and lets range through each one
			for _, checkout := range checkouts {

				// log.Println(checkout)

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

						txAddr := transfer.Payload_RPC.Value(
							rpc.RPC_REPLYBACK_ADDRESS,
							rpc.DataString,
						).(string)

						result, err := rpc.NewAddress(txAddr)
						if err != nil {
							return err
						}

						// as well as its destination port
						txPort := transfer.DestinationPort
						if txPort != addrPort ||
							txComment != addrComment ||
							txAddr != result.String() {
							continue
						}

						switch addrPort {
						case config.UserRegistrationPort:
							order := new(
								models.JSON_User_Order,
							)
							/*
								// this feel redundant
									// we already had "the order"
										// now we are having to make a new one
											order := I see what the problem is here,
											we are trying to rebuild the order because
											there is no order...

											we would need to be recording the details
											of the item to be able to retreive it.
											What this means is that we are able to handle
											the item like we do for the user; who is recorded
											in the rpc_payload (which is broken right now btw)

											anyway, I think that in order for this to work...
											you will have to record the item to the db and then
											you are going to have to do a "hasPaid" check on
											the function and then you are going to always
											be dealing with content that will be abused.

											The more ideal way for this would be to create a
											credits system... but then you are having to
											remember stuff, and that's not ideal.

											The thing that you are trying to do is make
											it so that the user is always having to interact
											with sending the server DERO; but I think that
											this model fails for 3 reasons:
												1. I think that a customer having to load^2
												more DERO into your website is good for you,
												but it isn't really all that good for them.
												They want to be able to be pulling deri out
												dayly, and not the other way around...
												2. you will need to need a way for them to
												make money on the platform. That will mean
												that you will need to host the data and make
												the items pay per view.
											So listings should be free, but the revenue will
											need to be made so that they are renting content.
							*/
							order.Name = txComment
							order.Wallet = result.String()
							if err := controllers.CreateUserRecord(order); err != nil {
								return err
							}
							fmt.Println("user created")
							user, err := controllers.GetUserByName(order.Name)
							if err != nil {
								return err
							}
							fmt.Println(user)
						}

					}
				}
			}
			time.Sleep(1 * time.Second)

		}
	}()
	return nil
}
