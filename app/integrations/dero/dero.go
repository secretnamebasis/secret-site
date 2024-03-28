package dero

import (
	"encoding/base64"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/secretnamebasis/secret-site/app/exports"
	"github.com/ybbus/jsonrpc"
)

func Call(
	endpoint string,
	object interface{},
	method string,
) error {
	// Create options for the JSON-RPC client
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " +
				base64.StdEncoding.EncodeToString(
					[]byte(
						config.Env("DERO_WALLET_USER")+
							":"+
							config.Env("DERO_WALLET_PASS"),
					),
				),
		},
	}
	rpcClient := jsonrpc.NewClientWithOpts(
		endpoint,
		opts,
	)

	err := rpcClient.CallFor(&object, method)

	return err
}

func Address() error {
	endpoint := "http://" +
		config.Env("DERO_SERVER_IP") +
		":" +
		config.Env("DERO_WALLET_PORT") +
		"/json_rpc"

	err := Call(
		endpoint,
		&exports.DeroAddressResult,
		"GetAddress",
	)

	if err != nil {
		return err
	}

	exports.DeroAddress,
		err = rpc.NewAddress(
		exports.DeroAddressResult.Address,
	)

	if err != nil {
		return err
	}

	return nil
}
