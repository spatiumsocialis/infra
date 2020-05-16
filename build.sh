docker build -t dependencies -f ./dependencies.Dockerfile .
docker-compose build
echo "Services built!"