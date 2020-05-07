# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=cmd/app
TEST_DIR=test
BINARY_NAME=app.out
SRC_NAME=main.go
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test build
build: 
	$(GOBUILD) -o ./$(BUILD_DIR)/${BINARY_NAME} -v ./$(BUILD_DIR)/${SRC_NAME}
test: 
	$(GOTEST) -v ./$(TEST_DIR)
.PHONY: test
clean: 
	$(GOCLEAN)
	rm -f ./$(BUILD_DIR)/${BINARY_NAME}
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
deps:
	$(GOGET) mod download