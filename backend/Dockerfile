# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and postgresql-client for pg_isready
RUN apk --no-cache add ca-certificates postgresql-client bash

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/migrate .

# Copy entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Copy migration files (needed for embedded migrations at runtime)
COPY --from=builder /app/internal/migrations ./internal/migrations

# Copy source code for runtime (needed for embedded migrations)
COPY --from=builder /app/internal ./internal

# Expose port (if you add an API server later)
EXPOSE 8080

# Use entrypoint script
ENTRYPOINT ["./entrypoint.sh"]
CMD ["up"]
