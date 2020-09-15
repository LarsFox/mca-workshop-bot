-include .env
export

default:
	@printf "Start with 'make init'"

init:
	@cp .env.example .env
	@printf "Set up the envs in '.env'!"

docker-build:
	@docker build -t mca_workshop_bot .

docker-run:
	@docker-compose up -d

lint:
	@golangci-lint run

run:
	@go run cmd/bot/main.go

mock:
	@go run cmd/model/main.go
