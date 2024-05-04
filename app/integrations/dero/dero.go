package dero

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	c "github.com/secretnamebasis/secret-site/app/config"
	"github.com/ybbus/jsonrpc"
)

const prefix = "DERO."
const user = "DERO_WALLET_USER"
const pass = "DERO_WALLET_PASS"
const DERO_SCID_STRING = "0000000000000000000000000000000000000000000000000000000000000000"

// CallRPC is a generic function to make JSON-RPC calls to either the DERO wallet or node.
func CallRPC(
	endpoint string,
	object interface{},
	method string,
	params ...interface{},
) error {

	// For DERO Node calls
	if strings.Contains(method, prefix) {
		rpcClient := jsonrpc.NewClient(endpoint)
		if len(params) > 0 {
			return rpcClient.CallFor(
				&object,
				method,
				params[0],
			)
		}
		return rpcClient.CallFor(
			&object,
			method,
		)
	}

	// For DERO Wallet calls
	endpointAuth := c.Env(c.EnvPath, user) + ":" + c.Env(c.EnvPath, pass)
	encodedEndpointAuth := base64.StdEncoding.EncodeToString([]byte(endpointAuth))

	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedEndpointAuth,
		},
	}

	rpcClient := jsonrpc.NewClientWithOpts(
		endpoint,
		opts,
	)

	if len(params) > 0 {
		return rpcClient.CallFor(
			&object,
			method,
			params[0],
		)
	}

	return rpcClient.CallFor(
		&object,
		method,
	)
}

// GetWalletAddress fetches the DERO wallet address.
func GetWalletAddress(endpoint string) (*rpc.Address, error) {
	// params := map[string]interface{}{}
	method := "GetAddress"
	err := CallRPC(endpoint, &c.ServerWallet, method)
	if err != nil {
		return nil, err
	}

	return rpc.NewAddress(
		c.ServerWallet.Address,
	)
}

func GetWalletTransfers(endpoint string) (*rpc.Get_Transfers_Result, error) {
	method := "GetTransfers"
	params := rpc.Get_Transfers_Params{}
	var response rpc.Get_Transfers_Result
	err := CallRPC(
		endpoint,
		&response,
		method,
		params,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetEncryptedBalance fetches the encrypted balance for the given address.
func GetEncryptedBalance(endpoint, address string) (*rpc.GetEncryptedBalance_Result, error) {
	method := prefix + "GetEncryptedBalance"
	params := rpc.GetEncryptedBalance_Params{
		Address:    address,
		TopoHeight: -1,
	}
	var response rpc.GetEncryptedBalance_Result
	err := CallRPC(
		endpoint,
		&response,
		method,
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
	method := prefix + "GetSC"
	err := CallRPC(
		endpoint,
		&response,
		method,
		rpc.GetSC_Params{
			SCID:      scid,
			Code:      true,
			Variables: true,
		},
	)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

func Comment(endpoint, comment, destination string) (rpc.Transfer_Result, error) {

	// and a pencil
	object := rpc.Transfer_Result{}
	// from the chart of accounts
	// turn to the leaf called "transfer"
	method := "transfer"
	transfer := rpc.Transfer{
		//
		SCID:        crypto.ZEROHASH,
		Destination: destination,
		Amount:      1, // we want them to keep one,
		Payload_RPC: rpc.Arguments{
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
		CallRPC(
			endpoint,
			&object,
			method,
			params,
		)
}

func MintContract(
	endpoint,
	contract,
	destination string,
) (rpc.Transfer_Result, error) {

	t := rpc.Transfer{
		Destination: destination,
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
	if err := CallRPC(
		endpoint,
		&obj,
		method,
		params,
	); err != nil {
		return rpc.Transfer_Result{}, err
	}

	return obj, nil
}

func MakeIntegratedAddress(
	comment string,
	price uint64,
	expiry time.Time,
) (rpc.Make_Integrated_Address_Result, error) {

	params := rpc.Make_Integrated_Address_Params{
		Address: c.ServerWallet.Address,
		Payload_RPC: rpc.Arguments{
			rpc.Argument{
				Name:     rpc.RPC_COMMENT,
				DataType: rpc.DataString,
				Value:    comment,
			},
			rpc.Argument{
				Name:     rpc.RPC_VALUE_TRANSFER,
				DataType: rpc.DataUint64,
				Value:    price,
			},
			rpc.Argument{
				Name:     rpc.RPC_EXPIRY,
				DataType: rpc.DataTime,
				Value:    expiry,
			},
		},
	}
	var result rpc.Make_Integrated_Address_Result
	method := "MakeIntegratedAddress"
	if err := CallRPC(
		c.WalletEndpoint,
		&result,
		method,
		params,
	); err != nil {
		return rpc.Make_Integrated_Address_Result{}, err
	}
	return result, nil
}
