#!/bin/sh

docker-compose run --rm start_dependencies
docker-compose up -d
echo "Services up and running!"
echo "Traefik dashboard available at ${DOCKERHOST}:8080"
echo "Services available at ${DOCKERHOST}:80"