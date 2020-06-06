#!/bin/bash
# Arguments
# $1: Project root dir
# $2: Build deploy dir
# $3: Specific service (optional)
docker-compose -f $2/docker-compose.yml pull $3
