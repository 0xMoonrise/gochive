FROM golang:1.26.4-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y unzip sudo

COPY setup.sh .
RUN ./setup.sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./cmd/gochive

FROM debian:bookworm-slim AS runtime

WORKDIR /app

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app ./
COPY --from=builder /usr/local/lib/libpdfium.so /usr/local/lib/

COPY --from=builder /app/static ./static

COPY --from=builder /opt/gochive/lib/pdfjs/web /opt/gochive/lib/pdfjs/web
COPY --from=builder /opt/gochive/lib/pdfjs/build /opt/gochive/lib/pdfjs/build

RUN ldconfig

CMD ["./app"]
