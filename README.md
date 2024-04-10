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
This site assumes TLS certification has been done in advance
```go
"/etc/letsencrypt/live/"+config.Domain+"/cert.pem",
"/etc/letsencrypt/live/"+config.Domain+"/privkey.pem"
```

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
Use can parse flags to customize your development environment. 
```sh
go run . -env=dev -port=3000 -db=./app/database/
```

## Testing
When you `run_integration_test.sh`, you will find a timesstamped build in `./builds/` and logs in `./logs/`
