include .env

ENTRY=./cmd/gochive
TARGET=gochive
ROOT=/opt/gochive
BACKUP=/mnt/usb/backups/gochive
CGO_ENABLED=1
MIGRATIONS_DIR=db/migrations

.PHONY: db 

all:
	go run $(ENTRY)

build:
	go build $(LDFLAGS) -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

backup:
	rsync -avu --stats "$(ROOT)/" "$(BACKUP)/"
	@echo "Done."

restore:
	rsync -avu --stats "$(BACKUP)/" "$(ROOT)/"
	@echo "Done."

test:
	go test ./...

db:
	sqlite3  -header -column $(ROOT)/gochive.db

sqlc:
	sqlc generate -f db/sqlc.yml
tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(ROOT)/gochive.db goose -dir $(MIGRATIONS_DIR) up

migrate-down:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(ROOT)/gochive.db goose -dir $(MIGRATIONS_DIR) down

migrate-status:
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=$(ROOT)/gochive.db goose -dir $(MIGRATIONS_DIR) status
