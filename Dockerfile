# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o moortown-house-backend main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/moortown-house-backend ./moortown-house-backend
COPY .env .env
EXPOSE 8080
CMD ["./moortown-house-backend"]

