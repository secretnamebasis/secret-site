# SECRET-SITE
## Intro
This project is aimed at creating a very hardy and robust [Go Fiber](https://gofiber.io/) site that is integrated with [DERO](https://dero.io)

## Arch
The hefitest part of this project lies squarely on `GoFiber`

The database is `bbolt`, an in app key/value store. 

The coolest feature included is the integration with `DERO`

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
- first-run script would be kind of cool
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

Use the [`dot.env.sample`](https://github.com/secretnamebasis/secret-site/blob/main/dot.env.sample) file to create `.env` files:

Duplicate the `.env` parmeters for the following directories:

#### project directory: `./.env` 

Fill out the `.env` variables, then copy and paste into project directory:
```sh
cat <<'CONFIG' > ./.env
# P U B L I C
## APP 
#
DOMAIN="example.site"

# P R I V A T E
## APP
#
SECRET="secretWords&Numbers2"

## DERO

### IP
#
DERO_NODE_IP="127.0.0.1"
DERO_WALLET_IP="127.0.0.1"

### PORT
#
DERO_NODE_PORT="10102"
DERO_WALLET_PORT="10103"

### AUTH
# 
DERO_WALLET_USER="secret"
DERO_WALLET_PASS="pass"

CONFIG

```

Or copy the template and fill out `.env` variables:
```sh
cp dot.env.sample .env 
nano .env
```
#### test directory: `./test/.env`
Once you have made your `.env` for the project, create `./test/` dependant variables ; our example assumes that, at first launch, both production (`prod`) and testing (`test`) instnaces are the same.

```sh
cp .env ./test/.env 
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
When you `run_integration_test.sh`, you will find a timesstamped builds in `./build/` and logs in `./log/`.

Alternatively, if you would like to test only the API:
```sh
bin/test
```
