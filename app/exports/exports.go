package exports

import (
	"github.com/deroproject/derohe/rpc"
)

const (
	APP_NAME    = "secret-site"
	DEV_ADDRESS = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
)

var (
	Env               string
	Port              int
	ProjectDir        = "./"
	EnvPath           = ProjectDir + ".env"
	DatabaseDir       string
	DeroAddress       *rpc.Address
	DeroAddressResult rpc.GetAddress_Result
)
