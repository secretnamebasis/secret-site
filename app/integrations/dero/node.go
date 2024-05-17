package dero

import "github.com/deroproject/derohe/rpc"

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
