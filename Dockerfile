# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -o main -ldflags="-w -s" .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy .env file if it exists (for development)
COPY --from=builder /app/.env ./.env 2>/dev/null || true

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 