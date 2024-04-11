package dero

import (
	"encoding/base64"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-site/app/config"
	"github.com/ybbus/jsonrpc"
)

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

// GetEncryptedBalanceResponse represents the JSON-RPC response for encrypted balance.
type GetEncryptedBalanceResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		SCID         string `json:"scid"`
		Data         string `json:"data"`
		Registration int    `json:"registration"`
		Bits         int    `json:"bits"`
		Height       int    `json:"height"`
		TopoHeight   int    `json:"topoheight"`
		BlockHash    string `json:"blockhash"`
		TreeHash     string `json:"treehash"`
		DHeight      int    `json:"dheight"`
		DTopoHeight  int    `json:"dtopoheight"`
		DTreeHash    string `json:"dtreehash"`
		Status       string `json:"status"`
	} `json:"result"`
}

// GetEncryptedBalance fetches the encrypted balance for the given address.
func GetEncryptedBalance(address string) (*GetEncryptedBalanceResponse, error) {

	params := map[string]interface{}{
		"address":    address,
		"topoheight": -1,
	}
	var response GetEncryptedBalanceResponse
	err := CallRPCNode(nodeEndpoint, &response, "DERO.GetEncryptedBalance", params)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
