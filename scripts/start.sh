# Arguments
# $1: Project root dir
# $2: path to dir containing Dockerfiles
# $3: environment {"dev", "prod"}
# $4: service (optional)

source ./scripts/dockerhost_mac_set.sh
docker-compose -f $2/docker-compose.yml run --rm start_dependencies
docker-compose -f $2/docker-compose.yml -f $2/docker-compose.$3.yml up -d $4
echo "Service(s) up and running!"
echo Jaeger tracing dashboard available at http://${DOCKERHOST}:16686
echo Traefik load balancer dashboard available at http://${DOCKERHOST}:8080
echo Services available at http://${DOCKERHOST}:80
