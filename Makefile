# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=cmd/app
BINARY_NAME=app.out
SRC_NAME=app.go
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test build
build: 
	$(GOBUILD) -o ./$(PACKAGE)$(BUILD_DIR)/${BINARY_NAME} -v ./$(PACKAGE)/$(BUILD_DIR)/${SRC_NAME}
test: 
	$(GOTEST) -v ./$(PACKAGE)...
.PHONY: test
clean: 
	$(GOCLEAN)
	rm -f ./$(PACKAGE)$(BUILD_DIR)/${BINARY_NAME}
run: build
	./$(PACKAGE)$(BUILD_DIR)/$(BINARY_NAME)
deps:
	$(GOGET) mod download