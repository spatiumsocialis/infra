# :mask: spatiumsocialis/infra

### Setup
1. Install Docker Desktop https://docs.docker.com/docker-for-mac/install/
2. Ask Matt for the .env file and Google service account JSON

## Commands
| Command                | Description                               |
|------------------------|-------------------------------------------|
| `make build`           | Build the services                        |
| `make start`           | Start the services                        |
| `make stop`            | Stop the services                         |
| `make test`            | Run the tests                             |
| `make coverage`        | Show HTML test coverage report            |
| `make dockerhost-mac`  | Output the local Docker IP                |
| `make token uid={uid}` | Generate a new access token for UID={uid} |
