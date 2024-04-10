# SECRET-SITE
## Intro
This project is aimed at creating a very hardy and robust Go Fiber site that is integrated with DERO
## Arch
The hefitest part of this project lies squarely on `GoFiber`

The coolest feature included is the integration with `DERO`

## Design
`secret-site` currently supports:
- bbolt database with encryption and decryption for item content (description, image)
    - Items { title: , content:, :description, :image }
    - Users { user: , wallet: }

## Roadmap
- ~~encrypt & decrypt content stored in database~~
- User authentication
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

## Testing
When you `run_integration_test.sh`, you will find a timesstamped builds in `./builds/` and logs in `./logs/`
