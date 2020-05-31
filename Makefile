# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=./build
BUILD_DEPLOY_DIR=$(BUILD_DIR)/deploy
BUILD_PACKAGE_DIR=$(BUILD_DIR)/package
SERVICE_DOCKERFILE=$(BUILD_PACKAGE_DIR)/service.Dockerfile
SERVICES=circle proximity scoring
EXECUTABLES=app consumer
BINARY_NAME_APP=app.out
BINARY_NAME_CONSUMER=consumer.out
BINARY_UNIX=$(BINARY_NAME)_unix
GOOGLE_GCR_HOSTNAME=gcr.io
GOOGLE_PROJECT_ID=spatiumsocialis

all: deps test build
# TODO: Clean this mess up
test: 
	$(GOTEST) ./$(package)... $(ARGS)
.PHONY: test
build-token:
	$(GOBUILD) -o ./tools/tokengen/cmd/tokengen/tokengen.out ./tools/tokengen/cmd/tokengen
token:
	./tools/tokengen/cmd/tokengen/tokengen.out -u $(uid)
push-deps:
	docker push ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest
push:
	docker-compose -f ${BUILD_DEPLOY_DIR}/docker-compose.yml -f ${BUILD_DEPLOY_DIR}/docker-compose.build.yml push ${service}
pull:
	docker-compose pull ${service}
build-deps:
	docker build -t ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest -f ${BUILD_PACKAGE_DIR}/deps.Dockerfile ${PWD}
start:
	docker-compose -f ${BUILD_DEPLOY_DIR}/docker-compose.yml run --rm start_dependencies
	docker-compose -f ${BUILD_DEPLOY_DIR}/docker-compose.yml -f ${BUILD_DEPLOY_DIR}/docker-compose.${env}.yml up -d ${service}
	@echo "Service(s) up and running!"
	@echo Jaeger tracing dashboard available at http://${DOCKERHOST}:16686
	@echo Traefik load balancer dashboard available at http://${DOCKERHOST}:8080
	@echo Services available at http://${DOCKERHOST}:80
build: build-deps
	PROJECT_ROOT=${PWD} SERVICE_DOCKERFILE=${SERVICE_DOCKERFILE} docker-compose -f ${BUILD_DEPLOY_DIR}/docker-compose.yml -f ${BUILD_DEPLOY_DIR}/docker-compose.build.yml build ${service}
	@echo "Service(s) built!"

stop:
	docker-compose ${BUILD_DEPLOY_DIR}/docker-compose.yml down ${service}
	@echo Services torn down
deploy:
	sh ./scripts/deploy.sh
ssh:
	gcloud beta compute ssh --zone "us-central1-a" "spatium-prod" --project "spatiumsocialis"
dockerhost:
	sh ./scripts/dockerhost.sh
