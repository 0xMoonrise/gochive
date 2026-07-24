# Gochive

Gochive is a personal project to store, back up, and centralize PDF documents, books, and articles. The main goal of this project is to practice and improve my development skills.

## Features

- Upload, search, favorite, edit, and delete files
- PDF viewer (via pdf.js) with thumbnail generation
- Markdown viewer with syntax highlighting, Mermaid diagrams, and MathJax equations
- Pluggable storage backend: local filesystem or S3-compatible storage
- SQLite database, managed with goose migrations
- CLI (built on Cobra) for running the server, backups, thumbnail regeneration, importing files from a URL, and inspecting the active configuration

## Requirements

- SQLite3
- A C toolchain (gcc, pkg-config) for building against pdfium

## Setup

`setup.sh` installs the native and frontend dependencies that gochive needs at runtime: pdfium, pdf.js, highlight.js, Mermaid, MathJax, marked, and DOMPurify. It downloads pinned versions of each and places them under `/opt/gochive/lib`.

```sh
./setup.sh
```

## Configuration

Gochive is configured via a TOML file. The configuration file is resolved in the following order. The first candidate found is used.

1. The path set in the `GOCHIVE_CONFIG` environment variable, if present.
2. `/opt/gochive/config.toml`.
3. `./config.toml`, relative to the current working directory.
4. `~/.config/gochive/config.toml`.

If none of these paths resolve to an existing file, startup fails with an explicit error.

A documented `config.example.toml` ships in the repo and inside the Docker image as a reference. Copy it to `config.toml` and adjust it for your setup. Never commit your real `config.toml`. It is excluded via `.gitignore`.

```toml
# Basic configuration
# Storage mode: 1 = filesystem, 2 = s3
mode = 1

# Address and port the server listens on
host = "0.0.0.0"
port = "8080"

# Directory where the SQLite database lives ($DATA/gochive.db)
data = "/opt/gochive/"

# Filesystem storage (required if mode = 1)
[fs]
root = "/opt/gochive/"

# S3 storage (required if mode = 2, omit otherwise)
# [s3]
# bucket = "gochive"
# access_key = "..."
# secret_key = "..."
# s3_endpoint = "https://s3.example.com"
# region = "us-east-1"
```

| Field | Description |
|---|---|
| `mode` | Storage backend selector: `1` for local filesystem, `2` for S3-compatible storage. |
| `host` / `port` | Address and port the HTTP server binds to. |
| `data` | Where the SQLite database file lives (`$data/gochive.db`). |
| `fs.root` | Where the filesystem storage backend keeps uploaded files (only used when `mode = 1`). |
| `s3.*` | S3-compatible client settings, only used when `mode = 2`. |

> `fs.root` and `data` are separate fields because storage and database location are decoupled by design. In most setups they point to the same path, but they can diverge, for example keeping the database on faster local disk while archived files live elsewhere.

### S3 secrets

`s3.access_key` and `s3.secret_key` can be set in `config.toml`, but for production deployments it is recommended to keep them out of the file entirely and provide them via environment variables instead, which override the values from the file if present:

```bash
S3_ACCESS_KEY=...
S3_SECRET_KEY=...
```

This keeps credentials out of any file that persists on disk or gets baked into a container image, and makes rotating a key a matter of updating an env var rather than editing a config file on the server.

Run `gochive status` at any time to print the resolved configuration (mode, paths, and the relevant S3/filesystem settings) without starting the server.

## Running with Docker

Build the image:

```sh
docker build -t gochive:1.0 .
```

The image only ships `config.example.toml` for reference. Your real `config.toml` is never baked into the image and must be mounted at runtime:

```sh
docker run -d --name gochive \
  --mount type=bind,source="$(pwd)/config.example.toml",target=/opt/gochive/config.toml,readonly \
  -v /opt/gochive/:/opt/gochive/ \
# -e S3_ACCESS_KEY=... \
# -e S3_SECRET_KEY=... \
  -p 8080:8080 \
  gochive:1.0
```

(Omit the `S3_*` env vars entirely if you are running in `mode = 1`.)

The container starts the server automatically (`gochive server`).

## Database

Migrations live in `internal/core/db/migrations` and run against a SQLite file at `$data/gochive.db`.

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

## CLI

Everything is handled through a single binary, `gochive` (built on Cobra), which bundles both the HTTP server and the maintenance commands:

```sh
gochive server                  # start the HTTP server on $host:$port
gochive status                  # print the active configuration
gochive backup -p <path>        # back up storage objects and database
gochive restore -p <path>       # restore storage objects and database
gochive generate [id]           # regenerate a thumbnail by id, or all of them if omitted
gochive upload_archive <url>    # download a file from a URL and add it to the archive
gochive version                 # print the binary version
```
