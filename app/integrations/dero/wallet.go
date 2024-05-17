package dero

import (
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	c "github.com/secretnamebasis/secret-site/app/config"
)

func Transfer(endpoint string, params rpc.Transfer_Params) (rpc.Transfer_Result, error) {
	method := "transfer"
	result := rpc.Transfer_Result{}
	return result,
		CallRPC(
			endpoint,
			&result,
			method,
			params,
		)
}

// GetWalletAddress fetches the DERO wallet address.
func GetWalletAddress(endpoint string) (*rpc.Address, error) {
	// params := map[string]interface{}{}
	method := "GetAddress"
	var response rpc.GetAddress_Result
	err := CallRPC(endpoint, &response, method)
	if err != nil {
		return nil, err
	}

	return rpc.NewAddress(
		response.Address,
	)
}

func GetWalletTransfers(
	endpoint string,
	params rpc.Get_Transfers_Params,
) (*rpc.Get_Transfers_Result, error) {
	method := "GetTransfers"
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

func Comment(
	endpoint,
	comment,
	destination string,
) (rpc.Transfer_Result, error) {

	// and a pencil
	result := rpc.Transfer_Result{}
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

	return result,
		CallRPC(
			endpoint,
			&result,
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
	params rpc.Make_Integrated_Address_Params,
) (rpc.Make_Integrated_Address_Result, error) {
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

func SplitIntegratedAddress(
	endpoint string,
	params rpc.Split_Integrated_Address_Params,
) (rpc.Split_Integrated_Address_Result, error) {
	var result rpc.Split_Integrated_Address_Result
	method := "SplitIntegratedAddress"
	if err := CallRPC(
		endpoint,
		&result,
		method,
		params,
	); err != nil {
		return rpc.Split_Integrated_Address_Result{}, err
	}
	return result, nil
}
