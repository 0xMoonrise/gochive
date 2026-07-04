# Gochive

Gochive is a personal project to store, back up, and centralize PDF documents, books, and articles. The main goal of this project is to practice and improve my development skills.

## Features

- Upload, search, favorite, edit, and delete files
- PDF viewer (via pdf.js) with thumbnail generation
- Markdown viewer with syntax highlighting, Mermaid diagrams, and MathJax equations
- Pluggable storage backend: local filesystem or S3-compatible storage
- SQLite database, managed with goose migrations

## Requirements

- Go 1.24+
- SQLite3
- A C toolchain (gcc, pkg-config) for building against pdfium

## Setup

`setup.sh` installs the native and frontend dependencies that gochive needs at runtime: pdfium, pdf.js, highlight.js, Mermaid, and MathJax. It downloads pinned versions of each and places them under `/opt/gochive/lib`.

```sh
./setup.sh
```

This step is already handled inside the Dockerfile if you build the image instead of running gochive directly on the host.

## Running with Docker

Build the image:
```sh
docker build -t gochive:1.0 .
```

Run the container:
```sh
docker run --name gochive \
  -v /opt/gochive/:/opt/gochive/ \
  --env-file=.env \
  -p 8080:8080 \
  --rm gochive:1.0
```

## Configuration

`.env`
```bash
HOST=0.0.0.0
PORT=8080
MODE=1
```

`MODE` selects the storage backend: `1` for local filesystem, `2` for S3.

S3 client environment variables (required only when `MODE=2`):
```bash
BUCKET=
ACCESS_KEY=
SECRET_KEY=
S3_ENDPOINT=
REGION=
```

## Database

Migrations live in `db/migrations` and run against a SQLite file at `/opt/gochive/gochive.db`.

```sh
make migrate-up      # apply migrations
make migrate-down    # roll back the last migration
make migrate-status   # check current migration status
make db               # open an interactive SQLite shell
```

## Backups

```sh
make backup     # rsync /opt/gochive to the configured backup path
make restore    # restore from the backup path
```

## CLI

Besides running as a server, the binary accepts a subcommand:

```sh
./gochive normalize
```

`normalize` migrates files stored under their old filename-based keys to the current id-based storage layout. Running the binary with no arguments starts the HTTP server.