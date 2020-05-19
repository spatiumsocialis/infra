# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR_APP=cmd/app
BUILD_DIR_CONSUMER=cmd/consumer
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
	$(GOTEST) -v ./$(PACKAGE)... $(ARGS)
.PHONY: test
clean-circle:
	$(GOCLEAN)
	rm -f ./circle/$(BUILD_DIR)/${BINARY_NAME}
clean-proximity:
	$(GOCLEAN)
	rm -f ./proximity/$(BUILD_DIR)/${BINARY_NAME}
clean-scoring:
	$(GOCLEAN)
	rm -f ./scoring/$(BUILD_DIR)/${BINARY_NAME}
clean: clean-circle clean-proximity clean-scoring
token:
	./auth/cmd/tokengen/tokengen.out -u $(uid)
run:
	./$(PACKAGE)$(BUILD_DIR)/$(EXEC) $(ARGS)
push-deps:
	docker push ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest
push:
	docker-compose -f docker-compose.yml -f docker-compose.build.yml push ${service}
pull:
	docker-compose pull ${service}
build-deps:
	docker build -t ${GOOGLE_GCR_HOSTNAME}/${GOOGLE_PROJECT_ID}/deps:latest -f ./deps.Dockerfile .
start:
	docker-compose run --rm start_dependencies
	docker-compose up -d ${service}
	@echo "Service(s) up and running!"
	@echo Jaeger tracing dashboard available at http://${DOCKERHOST}:16686
	@echo Traefik load balancer dashboard available at http://${DOCKERHOST}:8080
	@echo Services available at http://${DOCKERHOST}:80
build: build-deps
	docker-compose -f docker-compose.yml -f docker-compose.build.yml build ${service}
	@echo "Service(s) built!"

stop:
	docker-compose down ${service}
	@echo Services torn down