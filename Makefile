SHELL := /bin/sh

APP_NAME := ecommerce-shop
DB_URL ?= postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable

.PHONY: build run test docker-up docker-down migrate

build:
	GO111MODULE=on CGO_ENABLED=0 go build -o bin/app ./cmd/api

run:
	APP_ENV=development HTTP_ADDR=:8080 DATABASE_URL=$(DB_URL) JWT_SECRET=dev-secret go run ./cmd/api

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-unit:
	go test -v -short ./...

test-integration:
	go test -v -run Integration ./...

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

migrate:
	psql "$(DB_URL)" -f migrations/000_extensions.sql; \
	psql "$(DB_URL)" -f migrations/001_init.sql


