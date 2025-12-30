include .env

DATE := $(shell date '+%Y-%m-%d')

ENTRY=./cmd/gochive
TARGET=gochive
DOCKER_DB=gochive_db

.PHONY: db 

all:
	go run $(ENTRY)

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

backup:
	docker exec -it $(DOCKER_DB) pg_dump -U $(DB_USER) -d $(DB_NAME) | gzip > $(PATH_BACKUP)

restore:
	gunzip -c $(PATH_BACKUP) | docker exec -i $(DOCKER_DB) psql -U $(DB_USER) -h $(DB_HOST) -d $(DB_NAME)

test:
	go test ./...

db:
	docker exec -it $(DOCKER_DB) psql -U $(DB_USER) -d $(DB_NAME)

