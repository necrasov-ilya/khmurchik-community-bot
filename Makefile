.PHONY: run build test lint migrate-up migrate-down docker-up docker-down

run:
	go run ./cmd/bot/

build:
	go build -o bin/bot ./cmd/bot/

test:
	go test -v ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/khmurchik_bot?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/khmurchik_bot?sslmode=disable" down 1

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
