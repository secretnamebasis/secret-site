package dero

import (
	"encoding/base64"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/ybbus/jsonrpc"
)

const DERO_SCID_STRING = "0000000000000000000000000000000000000000000000000000000000000000"

var ()

// CallRPCNode is a generic function to make JSON-RPC calls to the DERO node.
func CallRPCNode(endpoint string, object interface{}, method string, params interface{}) error {
	rpcClient := jsonrpc.NewClient(endpoint)
	return rpcClient.CallFor(&object, method, params)
}

// CallRPCWallet is a generic function to make JSON-RPC calls to the DERO wallet.
func CallRPCWalletWithParams(endpoint string, object interface{}, method string, params interface{}) error {
	endpointAuth := config.Env(
		config.EnvPath,
		"DERO_WALLET_USER",
	) +
		":" +
		config.Env(
			config.EnvPath,
			"DERO_WALLET_PASS",
		)
	encodedEndpointAuth := base64.StdEncoding.EncodeToString([]byte(endpointAuth))
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedEndpointAuth,
		},
	}
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)

	return rpcClient.CallFor(
		&object,
		method,
		params,
	)
}

// CallRPCWallet is a generic function to make JSON-RPC calls to the DERO wallet.
func CallRPCWalletWithoutParams(endpoint string, object interface{}, method string) error {
	endpointAuth := config.Env(
		config.EnvPath,
		"DERO_WALLET_USER",
	) +
		":" +
		config.Env(
			config.EnvPath,
			"DERO_WALLET_PASS",
		)
	encodedEndpointAuth := base64.StdEncoding.EncodeToString([]byte(endpointAuth))
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedEndpointAuth,
		},
	}
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)

	return rpcClient.CallFor(
		&object,
		method,
	)
}

// GetWalletAddress fetches the DERO wallet address.
func GetWalletAddress(endpoint string) (*rpc.Address, error) {
	// params := map[string]interface{}{}
	err := CallRPCWalletWithoutParams(endpoint, &config.DeroAddressResult, "GetAddress")
	if err != nil {
		return nil, err
	}

	return rpc.NewAddress(
		config.DeroAddressResult.Address,
	)
}

// GetEncryptedBalance fetches the encrypted balance for the given address.
func GetEncryptedBalance(endpoint, address string) (*rpc.GetEncryptedBalance_Result, error) {

	params := rpc.GetEncryptedBalance_Params{
		Address:    address,
		TopoHeight: -1,
	}
	var response rpc.GetEncryptedBalance_Result
	err := CallRPCNode(
		endpoint,
		&response,
		"DERO.GetEncryptedBalance",
		params,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetSCID fetches the SCID for the given TXID.
func GetSCID(endpoint, scid string) (*rpc.GetSC_Result, error) {

	var response rpc.GetSC_Result

	err := CallRPCNode(
		endpoint,
		&response,
		"DERO.GetSC",
		rpc.GetSC_Params{
			SCID: scid,
			Code: true,
			// Variables: true,
		},
	)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

func Comment(endpoint, comment, destionation string) (rpc.Transfer_Result, error) {

	// and a pencil
	object := rpc.Transfer_Result{}
	// from the chart of accounts
	// turn to the leaf called "transfer"
	method := "transfer"
	transfer := rpc.Transfer{
		//
		SCID:        crypto.ZEROHASH,
		Destination: destionation,
		Amount:      1, // we want them to keep one,
		Payload_RPC: rpc.Arguments{
			// the first thing we want to do is establish
			// the habit of getting passwords by wallet.
			rpc.Argument{
				Name:     rpc.RPC_COMMENT,
				DataType: rpc.DataString,
				Value:    comment,
			},
		},
	}
	params := rpc.Transfer_Params{
		Transfers: []rpc.Transfer{
			transfer,
		},
	}

	return object,
		CallRPCWalletWithParams(
			endpoint,
			&object,
			method,
			params,
		)
}

func MintContract(endpoint, contract, destionation string) (rpc.Transfer_Result, error) {
	t := rpc.Transfer{
		Destination: destionation,
		Amount:      0,
	}
	arg := rpc.Argument{
		Name:     "entrypoint",
		DataType: "S",
		Value:    "InitializePrivate",
	}
	args := rpc.Arguments{arg}

	params := rpc.Transfer_Params{
		Transfers: []rpc.Transfer{t},
		SC_Code:   contract,
		SC_Value:  0,
		SC_RPC:    args,
		Ringsize:  2,
	}
	obj := rpc.Transfer_Result{}
	method := "transfer"
	if err := CallRPCWalletWithParams(
		endpoint,
		&obj,
		method,
		params,
	); err != nil {
		return rpc.Transfer_Result{}, err
	}

	return obj, nil
}
