# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=cmd/app
SERVICES=circle proximity scoring
BINARY_NAME=app.out
SRC_NAME=app.go
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test build
build-circle:
	$(GOBUILD) -o ./circle/$(BUILD_DIR)/${BINARY_NAME} -v ./circle/$(BUILD_DIR)/${SRC_NAME}
build-proximity:
	$(GOBUILD) -o ./proximity/$(BUILD_DIR)/${BINARY_NAME} -v ./proximity/$(BUILD_DIR)/${SRC_NAME}
build-scoring:
	$(GOBUILD) -o ./scoring/$(BUILD_DIR)/${BINARY_NAME} -v ./scoring/$(BUILD_DIR)/${SRC_NAME}
build: build-circle build-proximity build-scoring
test: 
	$(GOTEST) ./$(PACKAGE)...
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
run: build
	./$(PACKAGE)$(BUILD_DIR)/$(BINARY_NAME)
deps:
	$(GOGET) mod download