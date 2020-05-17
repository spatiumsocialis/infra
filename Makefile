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

all: deps test build
# TODO: Clean this mess up
build-circle-app:
	$(GOBUILD) -o ./circle/$(BUILD_DIR_APP)/${BINARY_NAME_APP} -v ./circle/$(BUILD_DIR_APP)
build-proximity-app:
	$(GOBUILD) -o ./proximity/$(BUILD_DIR_APP)/${BINARY_NAME_APP} -v ./proximity/$(BUILD_DIR_APP)
build-scoring-app:
	$(GOBUILD) -o ./scoring/$(BUILD_DIR_APP)/${BINARY_NAME_APP} -v ./scoring/$(BUILD_DIR_APP)
build-location-app:
	$(GOBUILD) -o ./location/$(BUILD_DIR_APP)/${BINARY_NAME_APP} -v ./location/$(BUILD_DIR_APP)
build-apps: build-circle-app build-proximity-app build-scoring-app
build-circle-consumer:
	$(GOBUILD) -o ./circle/$(BUILD_DIR_CONSUMER)/${BINARY_NAME_CONSUMER} -v ./circle/$(BUILD_DIR_CONSUMER)
build-proximity-consumer:
	$(GOBUILD) -o ./proximity/$(BUILD_DIR_CONSUMER)/${BINARY_NAME_CONSUMER} -v ./proximity/$(BUILD_DIR_CONSUMER)
build-scoring-consumer:
	$(GOBUILD) -o ./scoring/$(BUILD_DIR_CONSUMER)/${BINARY_NAME_CONSUMER} -v ./scoring/$(BUILD_DIR_CONSUMER)
build-consumers: build-circle-consumer build-proximity-consumer build-scoring-consumer
build-all-bins: build-apps build-consumers
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
run:
	./$(PACKAGE)$(BUILD_DIR)/$(EXEC) $(ARGS)
deps:
	$(GOGET) mod download
start:
	docker-compose run --rm start_dependencies
	docker-compose up -d
	@echo Services up and running!
	@echo Traefik dashboard available at ${DOCKERHOST}:8080
	@echo Services available at ${DOCKERHOST}:80
build:
	docker build -t dep -f ./dep.Dockerfile .
	docker-compose build
	@echo Services built!
stop:
	docker-compose down
	@echo Services torn down