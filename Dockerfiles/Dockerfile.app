FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git poppler-utils libwebp-tools

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/gochive

CMD ["./app"]
