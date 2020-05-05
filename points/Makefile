# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BUILD_DIR=cmd/app
SOURCE_DIR=cmd/app
SOURCE_NAME=main.go
BINARY_NAME=app.out
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test build
build: 
	$(GOBUILD) -o ./$(BUILD_DIR)/$(BINARY_NAME) -v ./$(SOURCE_DIR)/${SOURCE_NAME}
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN) ./$(BUILD_DIR)
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
get:
	$(GOGET) -d -v ./...
install: get
	$(GOINSTALL) -v ./...
deps:
	$(GOCMD) mod download