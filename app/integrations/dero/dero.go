package dero

import (
	"encoding/base64"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/ybbus/jsonrpc"
)

const DERO_SCID_STRING = "0000000000000000000000000000000000000000000000000000000000000000"

var (
	nodeEndpoint        = "http://" + config.Env("DERO_NODE_IP") + ":" + config.Env("DERO_NODE_PORT") + "/json_rpc"
	walletEndpoint      = "http://" + config.Env("DERO_WALLET_IP") + ":" + config.Env("DERO_WALLET_PORT") + "/json_rpc"
	endpointAuth        = config.Env("DERO_WALLET_USER") + ":" + config.Env("DERO_WALLET_PASS")
	encodedEndpointAuth = base64.StdEncoding.EncodeToString([]byte(endpointAuth))
)

// CallRPCNode is a generic function to make JSON-RPC calls to the DERO node.
func CallRPCNode(endpoint string, object interface{}, method string, params interface{}) error {
	rpcClient := jsonrpc.NewClient(endpoint)
	err := rpcClient.CallFor(&object, method, params)
	return err
}

// CallRPCWallet is a generic function to make JSON-RPC calls to the DERO wallet.
func CallRPCWalletWithParams(endpoint string, object interface{}, method string, params interface{}) error {
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedEndpointAuth,
		},
	}
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)
	err := rpcClient.CallFor(&object, method, params)
	return err
}

// CallRPCWallet is a generic function to make JSON-RPC calls to the DERO wallet.
func CallRPCWalletWithoutParams(endpoint string, object interface{}, method string) error {
	opts := &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"Authorization": "Basic " + encodedEndpointAuth,
		},
	}
	rpcClient := jsonrpc.NewClientWithOpts(endpoint, opts)
	err := rpcClient.CallFor(&object, method)
	return err
}

// GetWalletAddress fetches the DERO wallet address.
func GetWalletAddress() (*rpc.Address, error) {
	// params := map[string]interface{}{}
	err := CallRPCWalletWithoutParams(walletEndpoint, &config.DeroAddressResult, "GetAddress")
	if err != nil {
		return nil, err
	}
	address, err := rpc.NewAddress(config.DeroAddressResult.Address)
	return address, err
}

// GetEncryptedBalance fetches the encrypted balance for the given address.
func GetEncryptedBalance(address string) (*rpc.GetEncryptedBalance_Result, error) {

	params := map[string]interface{}{
		"address":    address,
		"topoheight": -1,
	}
	var response rpc.GetEncryptedBalance_Result
	err := CallRPCNode(nodeEndpoint, &response, "DERO.GetEncryptedBalance", params)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
func Comment(comment, destionation string) (rpc.Transfer_Result, error) {
	// grab up the General Journal
	endpoint := config.Env("DERO_NODE_IP")
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

	return object, CallRPCWalletWithParams(endpoint, &object, method, params)
}
