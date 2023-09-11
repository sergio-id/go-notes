.PHONY: sqlc-gen buf wire migrate_up migrate_down linter-golangci docker-up docker-down test-run

sqlc-gen:
	sqlc generate

buf:
	buf generate

wire:
	cd internal && \
    wire ./auth/app && \
    wire ./category/app && \
    wire ./note/app && \
    wire ./user/app

migrate_up:
	migrate -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' -path db/migrations up

migrate_down:
	migrate -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' -path db/migrations down

linter-golangci:
	golangci-lint run

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

test-run:
	go test ./internal/note/delivery/grpc
