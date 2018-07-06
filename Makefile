# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=core
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build: 
	$(GOBUILD) -v
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -v
	./$(BINARY_NAME)