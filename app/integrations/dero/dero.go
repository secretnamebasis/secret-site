package dero

import (
	"encoding/base64"
	"strings"

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
