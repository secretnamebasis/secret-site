# SECRET-SITE
## Intro
This project is aimed at creating a hardy and robust [Go Fiber](https://gofiber.io/) site that is integrated with [DERO](https://dero.io).
## Arch
The hefitest part of this project lies squarely on `GoFiber`.

The database is `bbolt`, an in app key/value store. 

The coolest feature included is the integration with `DERO`.
## Design
`secret-site` is a TLS encrypted website for hosting encrypted content. 

The following models are supported in the `bbolt` database with the accompanying features: 
- Item with`AES` encryption/decryption of `data:`
- User with wallet address validations for `DERO` network
## Roadmap
### DOCS
- API documentation 
### DB
- db encryption migrations
- routine backups
    - segmented backups (conserve storage)
#### ITEM
- ~~`AES` encrypted items~~
    - ~~`:description`~~
    - ~~`:image`~~
    - user authenticated, `AES` encrypted items
#### USER
- authentication
- signup
- login/logout
### EXTRAS
#### BACKEND
- ~~config script~~
- websocket connections with DERO wallets would be rad 
#### FRONTEND
- [`Gnomon`](https://github.com/civilware/Gnomon) search tools
- [`NFA`](https://github.com/civilware/artificer-nfa-standard) minting tools
## Install
### DERO wallet
As a pre requisite, running `secret-site` in production environments requires a DERO wallet instance:
- CLI:
    - Download the latest binaries of [`DERO`](https://github.com/deroproject/derohe/releases/latest/)
    - Restore, or create, a wallet file
    - Launch `dero-wallet-cli` with these suggested flags ; supposing of course that `derod` runs locally:
```sh
--rpc-bind=127.0.0.1:10103 \
--daemon-address=127.0.0.1:10102 \
--rpc-server \
--rpc-login="secret:pass"
```
- GUI: 
    - Download the latest [`ENGRAM`](https://github.com/DEROFDN/Engram/releases/latest/)
    - Restore, or create, wallet file
    - Configure module `Cyberdeck` for wallet `RPC` activation
    - Modify &/or collect `user` `pass` details
### Clone
Clone repo and change directories:
```sh
git clone https://github.com/secretnamebasis/secret-site.git
cd secret-site
```
## Config
It is assume that on first `config`, that production (`prod`), development (`dev`) and testing (`test`) are the same. 
### `.env`
Default values in [`dot.env.sample`](https://github.com/secretnamebasis/secret-site/blob/main/dot.env.sample) are used to set default values for the `.env` variables prior to running the `config`, which will write `env` files to the project directory `./`: 
- `.env` 
- `.env.dev`
- `.env.test` 
```sh
bin/config
```  
### SSL cert
This site [assumes SSL certification](https://github.com/secretnamebasis/secret-site/blob/cd559806442bad5553464d6fbee86966fec1aa3e/app/site.go#L41).
### Run
To run the application: 
```sh
go run .
``` 
or, if you prefer:  
```sh
go build . 
./secret-site
```
## Development/Testing
The `DERO` `simulator` runs in the background for all developemnt and testing environments.

### `dev`
Any `env` but `prod` runs app without TLS. Use parse flags to customize your development environment. 
```sh
go run . -env=dev -port=3000 -db=./app/database/
```
Or:
```sh
bin/dev
```
### `test`
When you `run_integration_test.sh`, you will find times-stamped builds in `./build/` and logs in `./log/`.

Alternatively, if you would like to test only the API:
```sh
bin/test
```
### Releases
We have also included a helpful `gh` script for deploying releases to `GitHub`
```sh
bin/release
```
