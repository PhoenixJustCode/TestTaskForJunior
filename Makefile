# Makefile для TestTaskForJun

SERVICE_NAME=test-task-app

# Переменные для подключения к БД
DB_HOST ?= db
DB_PORT ?= 5432
DB_USER ?= postgres
DB_NAME ?= bookTest_db
DB_PASSWORD ?= mysecretpassword
DB_SSLMODE ?= disable

.PHONY: build run migrate swag clean

build:
	docker-compose build $(SERVICE_NAME)

run:
	docker-compose up $(SERVICE_NAME)

migrate:
	migrate -path ./schema -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)" up

swag:
	swag init -g cmd/main.go

clean:
	docker-compose down -v
