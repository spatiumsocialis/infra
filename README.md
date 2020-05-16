# infra

## Quickstart guide
### Setup
1. Install Docker Desktop https://docs.docker.com/docker-for-mac/install/
2. From the project root, run the following to set the `DOCKERHOST` env variable
```
source dockerhost.sh
```
3. Create a .env file with the following variables
```
DB_PROVIDER=sqlite3
DB_CONNECTION_STRING=:memory:
PORT=8080
GOOGLE_API_KEY="your google api key"
GOOGLE_APPLICATION_CREDENTIALS=path/to/google/service/account.json
```
### Build
Run `make build` to build the services
### Start
Run `make start` to start the services
### Stop
Run `make stop` to stop the services
