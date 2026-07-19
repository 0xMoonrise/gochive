ENTRY=./cmd/gochive
TARGET=gochive
MIGRATIONS_DIR=./internal/core/db/migrations

.PHONY: db 

all:
	go run $(ENTRY)

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

test:
	go test ./...

db:
	sqlite3  -header -column $(DATA)/gochive.db

sqlc:
	sqlc generate -f ./internal/database/sqlc.yml

tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(DATA)/gochive.db goose -dir $(MIGRATIONS_DIR) up

migrate-down:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(DATA)/gochive.db goose -dir $(MIGRATIONS_DIR) down

migrate-status:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(DATA)/gochive.db goose -dir $(MIGRATIONS_DIR) status
