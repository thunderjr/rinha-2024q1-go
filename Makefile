#!make
include .env
export $(shell sed 's/=.*//' .env)

run:
	go run ./cmd/main.go

build:
	GOOS=linux GOARCH=arm64 go build -o bin/server ./cmd/main.go

compose:
	docker compose -f bin/docker-compose.yml up -d

admin:
	docker run --network rinha-2024-q1-go -e ME_CONFIG_MONGODB_URL=${MONGODB_URI} -p 8081:8081 mongo-express
