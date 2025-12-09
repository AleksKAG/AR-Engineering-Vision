.PHONY: build run compose up migrate

build:
	go build -o bin/ar cmd/api/main.go

run:
	./bin/ar

compose:
	docker compose -f docker/docker-compose.yml up --build

migrate:
	@echo "Run migrations with your favorite tool (golang-migrate)"
