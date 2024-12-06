#/bin/bash
# Use Bash instead of SH
export SHELL := /bin/bash

.DEFAULT_GOAL := controll

GOPATH := $(shell go env GOPATH)

SERVER_PATH := server/cmd/app
CLIENT_PATH := client/cmd/app

# Run the server
app-server:
	@echo "Server running..."
	@go run -race $(SERVER_PATH)/main.go

# Run the client
app-client:
	@echo "Client running..."
	@go run -race $(CLIENT_PATH)/main.go
