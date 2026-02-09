include .env

ENTRY=./cmd/gochive
TARGET=gochive

.PHONY: db 

all:
	go run $(ENTRY)

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

backup:
	rsync -avu --stats "$(APP_ROOT)/" "$(BACKUP_ROOT)/"
	@echo "Done."

restore:
	rsync -avu --ignore-existing --stats "$(BACKUP_ROOT)/" "$(APP_ROOT)/"
	@echo "Done."

test:
	go test ./...

db:
	sqlite3 $(RDIR)/gochive.db
