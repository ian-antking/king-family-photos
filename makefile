MAKEFILE_PATH = $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR = $(dir $(MAKEFILE_PATH))
BIN_DIR = $(CURRENT_DIR)/bin

.PHONY: build clean deploy-dev deploy-live

build: 
	cd $(CURRENT_DIR)/resizePhoto; env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o $(BIN_DIR)/resizePhoto main.go

clean:
	rm -rf ./bin

deploy-live: SHELL:=/bin/bash
deploy-live: clean build
	serverless deploy --verbose --stage live

deploy-dev: SHELL:=/bin/bash
deploy-dev: clean build
	serverless deploy --verbose --stage dev

teardown-dev: SHELL:=/bin/bash
teardown-dev:
	serverless remove --verbose --stage dev

teardown-live: SHELL:=/bin/bash
teardown-live:
	serverless remove --verbose --stage live