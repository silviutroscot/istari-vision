# Installation
- Go 
- Node
- NPM

# Run
- `source .env` to load the env file (either prod or dev)
- `make build && ./out/executable` to build and run the backend
- `make build-frontend` to build the frontend code
- `caddy run` to start the front end
- Note: the endpoint where the front end sends the request is the remote one, you can change it to the local one, but it is hardcoded for now
