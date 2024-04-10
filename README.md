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
- ~~encrypt & decrypt stores in database~~
    - ~~`:description`~~
    - ~~`:image` ~~
- authenticated, `AES` encrypted items

#### USER
- authentication
- signup
- login

### FRONTEND
- `Gnomon` search tools
- `NFA` minting tools

### EXTRAS
- first-run script would be kind of cool
- websocket connections with DERO wallets would be rad 

## Install

### env
Use the `dot.env.sample` file to create `.env` files for the following directories:
- root directory, `./.env` 
- test directory, `./test/.env`

### TLS cert
This site [assumes TLS certification](https://github.com/secretnamebasis/secret-site/blob/cd559806442bad5553464d6fbee86966fec1aa3e/app/site.go#L41) has been done in advance.

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
or,
```sh
bin/dev
```

## Testing
When you `run_integration_test.sh`, you will find a timesstamped builds in `./build/` and logs in `./log/`.

Alternatively, if you would like to just test the API:
```sh
bin/test
```
