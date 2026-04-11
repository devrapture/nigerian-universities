.PHONY: migrate-up migrate-down migrate-create dev

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

dev:
	air
