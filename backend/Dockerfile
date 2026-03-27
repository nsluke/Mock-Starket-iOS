# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /seed cmd/seed/main.go

# Dev stage (for hot reload with Docker Compose)
FROM golang:1.23-alpine AS dev

WORKDIR /app
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

COPY --from=builder /server /usr/local/bin/server
COPY --from=builder /migrate /usr/local/bin/migrate
COPY --from=builder /seed /usr/local/bin/seed
COPY migrations /migrations

EXPOSE 8080

CMD ["server"]
