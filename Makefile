# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

export DOCKER_TAG?=latest

# App parameters
APP_DIR=./cmd
BIN_DIR=./bin
APPS=$(wildcard $(APP_DIR)/*)
BINARY_NAMES=$(patsubst $(APP_DIR)/%, $(BIN_DIR)/%, $(APPS))
DOCKER_RUN_CMD=docker run --rm -it -p 8080:8080 $(DOCKER_IMAGE_NAME)

all: test build

build: $(BINARY_NAMES)

%: $(APP_DIR)/%
	$(GOBUILD) -o $(BIN_DIR)/$@ $<


test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

docker-build-%:
	docker build -t $*:$(DOCKER_TAG) --build-arg APP_NAME=$* .

docker-run-%:
	docker run --rm -it -p 8080:8080 $*:$(DOCKER_TAG)

docker-clean-%:
	docker rmi $*:$(DOCKER_TAG)

.PHONY: all build test clean run docker-build docker-run docker-clean

