# Arguments
# $1: Project root dir
# $2: path to service.Dockerfile
# $3: Build deploy dir
# $4: Service name

source .env
PROJECT_ROOT=$1 SERVICE_DOCKERFILE=$2 docker-compose -f $3/docker-compose.yml -f $3/docker-compose.build.yml build $4
