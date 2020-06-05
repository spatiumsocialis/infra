# Arguments
# $1: path to dir containing Dockerfiles
# $2: environment {"dev", "prod"}
# $3: service (optional)

source ./scripts/dockerhost.sh
docker-compose -f $1/docker-compose.yml run --rm start_dependencies
docker-compose -f $1/docker-compose.yml -f $1/docker-compose.$2.yml up -d $3
echo "Service(s) up and running!"
echo Jaeger tracing dashboard available at http://${DOCKERHOST}:16686
echo Traefik load balancer dashboard available at http://${DOCKERHOST}:8080
echo Services available at http://${DOCKERHOST}:80
