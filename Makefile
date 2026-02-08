include .env

DATE := $(shell date '+%Y-%m-%d')
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
	rsync -av --ignore-existing --stats "$(RDIR)/" "$(BK_DIR)/"
	@echo "Done."

restore:
	rsync -av --ignore-existing --stats "$(BK_DIR)/" "$(RDIR)/"
	@echo "Done."

test:
	go test ./...

db:
	sqlite3 $(RDIR)/gochive.db
