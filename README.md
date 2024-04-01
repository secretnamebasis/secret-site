# SECRET-SITE
## Intro
This project is aimed at creating a very hardy and robust Go Fiber site that is integrated with DERO
## Design
`secret-site` currently supports: 
- Items { title: , content: }
- Users { user: , wallet: }

## Roadmap
-  current `db.go` is a `bbolt` implmentation, an encrypted database would be preferred

## Install
Use the `dot.env.sample` file to create `.env` files for the following directories:
- root directory, `./.env` 
- test directory, `./test/.env`

Then `go run .` or `go build . && ./secret-site`
## Testing
When you `run_integration_test.sh`, you will find a build in `./builds/` and logs in `./logs/`
