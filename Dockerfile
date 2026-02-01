FROM golang:1.24-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y unzip

COPY setup.sh .
RUN ./setup.sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build \
    -ldflags="-extldflags '-Wl,-rpath,/usr/local/lib'" \
    -o app ./cmd/gochive

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app ./
COPY --from=builder /usr/local/lib/libpdfium.so /usr/local/lib/

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/db/migrations ./db/migrations

COPY --from=builder /opt/pdfjs/web /opt/pdfjs/web
COPY --from=builder /opt/pdfjs/build /opt/pdfjs/build

RUN ldconfig

CMD ["./app"]
