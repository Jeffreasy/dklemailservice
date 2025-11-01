# Build stage for production (CGO disabled for smaller, static binaries)
FROM golang:1.23-alpine3.19 AS builder-prod

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations and CGO disabled for production
RUN CGO_ENABLED=0 GOOS=linux go build -o main-prod -ldflags="-w -s" .

# Build stage for development/testing (CGO enabled for SQLite support)
FROM golang:1.23-alpine3.19 AS builder-dev

WORKDIR /app

# Install necessary build tools and SQLite dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled for testing
RUN CGO_ENABLED=1 GOOS=linux go build -o main-dev .

# Final stage
FROM alpine:3.19 AS runtime

WORKDIR /app

# Install ca-certificates for HTTPS and SQLite for development
RUN apk --no-cache add ca-certificates sqlite

# Copy the binaries from builders
COPY --from=builder-prod /app/main-prod ./main
COPY --from=builder-dev /app/main-dev ./main-dev

# Copy templates directory
COPY --from=builder-prod /app/templates ./templates

# Create a .env file with default values for database
RUN echo "DB_HOST=\${DB_HOST:-localhost}" > .env && \
    echo "DB_PORT=\${DB_PORT:-5432}" >> .env && \
    echo "DB_USER=\${DB_USER:-postgres}" >> .env && \
    echo "DB_PASSWORD=\${DB_PASSWORD:-}" >> .env && \
    echo "DB_NAME=\${DB_NAME:-dklemailservice}" >> .env && \
    echo "DB_SSL_MODE=\${DB_SSL_MODE:-disable}" >> .env

# Expose port
EXPOSE 8080

# Set environment variable to indicate which binary to use
# Use ENV APP_ENV=dev to use the development binary with CGO enabled
ENV APP_ENV=prod

# Command to run the executable based on environment
CMD if [ "$APP_ENV" = "dev" ]; then ./main-dev; else ./main; fi 