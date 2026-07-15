# Gochive

Gochive is a personal project to store, back up, and centralize PDF documents, books, and articles. The main goal of this project is to practice and improve my development skills.

## Features

- Upload, search, favorite, edit, and delete files
- PDF viewer (via pdf.js) with thumbnail generation
- Markdown viewer with syntax highlighting, Mermaid diagrams, and MathJax equations
- Pluggable storage backend: local filesystem or S3-compatible storage
- SQLite database, managed with goose migrations
- CLI (built on Cobra) for backups, thumbnail regeneration, importing files from a URL, and inspecting the active configuration

## Requirements

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
docker run -d --name gochive \
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
ROOT=/opt/gochive/
DATA=/opt/gochive/
BACKUP=/mnt/usb/backups/gochive/
```

| Variable | Description |
|---|---|
| `MODE` | Storage backend selector: `1` for local filesystem, `2` for S3-compatible storage. |
| `ROOT` | Where the filesystem storage backend keeps uploaded files (only used when `MODE=1`). |
| `DATA` | Where the SQLite database file lives (`$DATA/gochive.db`). |
| `BACKUP` | Destination used by the CLI's `backup`/`restore` commands. |

> `ROOT` and `DATA` are separate variables because storage and database location are decoupled by design — in most setups they'll point to the same path, but they can diverge (e.g. keeping the DB on faster local disk while archived files live elsewhere).

S3 client environment variables (required only when `MODE=2`):

```bash
BUCKET=
ACCESS_KEY=
SECRET_KEY=
S3_ENDPOINT=
REGION=
```

Run `gochive-cli status` at any time to print the resolved configuration (mode, paths, and the relevant S3/filesystem settings) without starting the server.

## Database

Migrations live in `internal/core/db/migrations` and run against a SQLite file at `$DATA/gochive.db`.

```sh
make migrate-up      # apply migrations
make migrate-down    # roll back the last migration
make migrate-status  # check current migration status
make db              # open an interactive SQLite shell
```

Query code is generated with sqlc from `internal/database/sqlc.yml`:

```sh
make sqlc
```

## Server

```sh
./gochive
```

Starts the HTTP server on `$HOST:$PORT`.

## CLI

A separate `gochive-cli` binary (built from `cmd/gochive-cli`) provides maintenance and import commands:

```sh
gochive-cli status                 # print the active configuration
gochive-cli backup                 # rsync data to the configured backup path
gochive-cli restore                # restore data from the backup path
gochive-cli generate [id]          # regenerate a thumbnail by id, or all of them if omitted
gochive-cli upload_archive <url>   # download a file from a URL and add it to the archive
gochive-cli version                # print the CLI version
```
