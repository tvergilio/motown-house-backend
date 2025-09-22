# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o web-service-gin main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web-service-gin ./web-service-gin
COPY .env .env
EXPOSE 8080
CMD ["./web-service-gin"]

