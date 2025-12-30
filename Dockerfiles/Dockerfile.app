FROM golang:1.24.4-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    git wget ca-certificates pkg-config build-essential \
 && rm -rf /var/lib/apt/lists/*

COPY Dockerfiles/setup.sh ./
RUN chmod +x setup.sh && ./setup.sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o app ./cmd/gochive


FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/local/lib/libpdfium.so /usr/local/lib/
COPY --from=builder /app/app .
COPY --from=builder /app/.env .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/db/migrations /app/db/migrations

ENV LD_LIBRARY_PATH=/usr/local/lib

CMD ["./app"]
