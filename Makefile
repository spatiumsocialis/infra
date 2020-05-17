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
	docker build -t deps -f ./deps.Dockerfile .
start:
	docker-compose run --rm start_dependencies
	docker-compose up -d
	@echo Services up and running!
	@echo Traefik dashboard available at ${DOCKERHOST}:8080
	@echo Services available at ${DOCKERHOST}:80
build: deps
	docker-compose build
	@echo Services built!
stop:
	docker-compose down
	@echo Services torn down