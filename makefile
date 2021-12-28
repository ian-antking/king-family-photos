MAKEFILE_PATH = $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR = $(dir $(MAKEFILE_PATH))
BIN_DIR = $(CURRENT_DIR)/bin

.PHONY: build clean deploy

build: 
	cd $(CURRENT_DIR)/resizePhoto; env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o $(BIN_DIR)/resizePhoto main.go

clean:
	rm -rf ./bin

deploy: SHELL:=/bin/bash
deploy: clean build
	serverless deploy --verbose

teardown: SHELL:=/bin/bash
teardown:
	serverless remove --verbose