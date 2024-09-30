include .env

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

build:
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -o out

run: build
	@./out

clean:
	@rm -rf out
	@go mod tidy

up:
	@goose -dir ${SCHEMA_DIR} ${DRIVER} ${CONN_STRING} up

down:
	@goose -dir ${SCHEMA_DIR} ${DRIVER} ${CONN_STRING} down

models:
	sqlc generate

evolve: up models

.PHONY: build run clean up down models evolve