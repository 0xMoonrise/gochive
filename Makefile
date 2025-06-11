DATE := $(shell date '+%Y-%m-%d')

ENTRY=./cmd/gochive
TARGET=gochive

BACKUP_NAME=archive_db.$(DATE).gz
PATH_BACKUP=/tmp/$(BACKUP_NAME)

all:
	go run $(ENTRY)

build:
	go build -o $(TARGET) $(ENTRY)

clean:
	rm $(TARGET)

backup:
	docker exec -it archive_db pg_dump -U $(DB_USER) -d $(DB_NAME) | gzip > $(PATH_BACKUP)

restore:
	gunzip -c $(PATH_BACKUP) | docker exec -i archive_db psql -U $(DB_USER) -h 127.0.0.1 -d $(DB_NAME)
