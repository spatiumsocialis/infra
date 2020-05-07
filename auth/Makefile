# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test
test: 
	$(GOTEST) -v ./...
get:
	$(GOGET) -d -v ./...
install: get
	$(GOINSTALL) -v ./...
deps:
	$(GOGET) mod download