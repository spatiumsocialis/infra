# infra

## Quickstart guide
1. Install Docker Desktop https://docs.docker.com/docker-for-mac/install/
2. From the project root, run the following to set the DOCKERHOST env variable
```
export DOCKERHOST=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)
```
3. Create a .env file with the following variables
```
DB_PROVIDER=sqlite3
DB_CONNECTION_STRING=:memory:
PORT=8080
GOOGLE_API_KEY="your google api key"
GOOGLE_APPLICATION_CREDENTIALS=path/to/google/service/account.json
```
4. From the project root, run `docker build -t dependencies -f ./dependencies.Dockerfile .` to build the dependencies image
5. Run `docker-compose run --rm start_dependencies` to start the Kafka and Zookeeper services
6. Run `docker-compose up` to build and run the rest of the services
