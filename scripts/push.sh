# Arguments
# $1: Project root
# $2: Build package directory
# $3: Build deploy directory
# $4: Specific service name (optional)

source .env
export PROJECT_ROOT=$1
export SERVICE_DOCKERFILE=$2/service.Dockerfile
docker-compose -f $3/docker-compose.yml -f $3/docker-compose.build.yml push $4
