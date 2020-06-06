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
SOURCE_ENV=source .env

ifeq ($(env),)
env := dev
endif

all: deps test build
# TODO: Clean this mess up
test: 
	$(SOURCE_ENV) && $(GOTEST) -coverprofile=/tmp/coverage.out ./pkg/... $(ARGS)
.PHONY: test
coverage:
	go tool cover -html=/tmp/coverage.out
build-token:
	$(GOBUILD) -o ./tools/tokengen/cmd/tokengen/tokengen.out ./tools/tokengen/cmd/tokengen
token:
	$(SOURCE_ENV) && ./tools/tokengen/cmd/tokengen/tokengen.out -u $(uid)
push-deps:
	$(SOURCE_ENV) && docker push ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest
push:
	sh ./scripts/push.sh ${PWD} ${BUILD_PACKAGE_DIR} ${BUILD_DEPLOY_DIR} ${service}
pull:
	sh ./scripts/pull.sh ${PWD} ${BUILD_DEPLOY_DIR} ${service}
build-deps:
	sh ./scripts/build-deps.sh ${GOOGLE_GCR_HOSTNAME} ${GOOGLE_PROJECT_ID} ${BUILD_PACKAGE_DIR} ${PWD}
start:
	sh ./scripts/start.sh ${PWD} ${BUILD_DEPLOY_DIR} ${env} ${service}
build: build-deps
	sh ./scripts/build.sh ${PWD} ${SERVICE_DOCKERFILE} ${BUILD_DEPLOY_DIR} ${service} 
stop:
	docker-compose ${BUILD_DEPLOY_DIR}/docker-compose.yml down ${service}
	@echo Services torn down
deploy:
	sh ./scripts/deploy.sh
ssh:
	gcloud beta compute ssh --zone "us-central1-a" "spatium-prod" --project "spatiumsocialis"
dockerhost:
	sh ./scripts/dockerhost.sh
