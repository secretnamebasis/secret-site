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
- Item: `{ id: , title: , content: :description, :image, :imageURL }`
    - `AES` encryption/decryption of `:description`, `:image`
- User `{ user: , wallet: , password: }`
    - validates wallet addresses with `DERO` network
## Roadmap
### DOCS
- API documentation 
### DB
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
### clone
Clone repo and change directories:
```sh
git clone https://github.com/secretnamebasis/secret-site.git
cd secret-site
```
### DERO wallet
`secret-site` need a DERO wallet instance:
- Download the latest [`ENGRAM`](https://github.com/DEROFDN/Engram/releases/latest/) wallet
- Restore, or create, wallet file
- Configure module `Cyberdeck` for wallet `RPC` activation
- Modify &/or collect `user` `pass` details
- Paste `user` `pass` values into `.env` file
### env
Default values in [`dot.env.sample`](https://github.com/secretnamebasis/secret-site/blob/main/dot.env.sample) are used to set default values for the `.env` variables prior to running the `config`, which will write `.env` to the project directory `./.` and the `./test/` direcorty. We assume that on first `config`, that insteance of production (`prod`), development (`dev`) and testing (`test`) are the same. 
```sh
bin/config
```  
### TLS cert
This site [assumes TLS certification](https://github.com/secretnamebasis/secret-site/blob/cd559806442bad5553464d6fbee86966fec1aa3e/app/site.go#L41).
### run
To run the application: 
```sh
go run .
``` 
or, if you prefer:  
```sh
go build . 
./secret-site
```
## Dev 
Any `env` but `prod` runs app without TLS. Use parse flags to customize your development environment. 
```sh
go run . -env=dev -port=3000 -db=./app/database/
```
Or:
```sh
bin/dev
```
## Testing
When you `run_integration_test.sh`, you will find times-stamped builds in `./build/` and logs in `./log/`.

Alternatively, if you would like to test only the API:
```sh
bin/test
```
