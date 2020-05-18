export DOCKERHOST=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)
docker-compose run --rm start_dependencies
docker-compose up -d
echo "Service(s) up and running!"
echo Jaeger tracing dashboard available at http://${DOCKERHOST}:16686
echo Traefik load balancer dashboard available at http://${DOCKERHOST}:8080
echo Services available at http://${DOCKERHOST}:80