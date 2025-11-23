FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o pr_service ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/pr_service .

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .env

CMD ["./pr_service"]
