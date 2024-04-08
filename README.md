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
- ~~current `db.go` is a `bbolt` implmentation, an encrypted database would be preferred~~
- User authentication; turned off for the moment
- first-run script would be kind of cool
- websocket connections with DERO wallets would be rad 

## Install
Use the `dot.env.sample` file to create `.env` files for the following directories:
- root directory, `./.env` 
- test directory, `./test/.env`

Then `go run .` or `go build . && ./secret-site`
## Testing
When you `run_integration_test.sh`, you will find a build in `./builds/` and logs in `./logs/`
