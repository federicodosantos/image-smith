include .env

run:
	@go run cmd/main.go 

migrate-up:
	@migrate -path db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up    

migrate-down:
	@migrate -path db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

migrate-fix:
	@migrate -path db/migrations -database "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose force 20241106063649

.PHONY: run migrate-up migrate-down migrate-fix
