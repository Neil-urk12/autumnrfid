# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN go install github.com/air-verse/air@latest
# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o autum ./cmd/web/main.go
# Dev stage for live reload with Air
FROM builder AS dev
WORKDIR /app
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
# Final stage
FROM alpine:3.21 AS final
RUN apk add --no-cache ca-certificates
WORKDIR /root/

# Copy binary and assets
COPY --from=builder /app/autum .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/.env .env

# Expose port
EXPOSE 8080

# Run
CMD ["./autum"]
