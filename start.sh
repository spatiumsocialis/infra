#!/bin/sh

docker-compose run --rm start_dependencies
docker-compose up -d
echo "Services up and running!"