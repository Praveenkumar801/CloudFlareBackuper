# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cloudflare-backuper .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1000 backuper && \
    adduser -D -u 1000 -G backuper backuper

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/cloudflare-backuper .

# Copy example config
COPY config.example.yml .

# Change ownership
RUN chown -R backuper:backuper /app

# Switch to non-root user
USER backuper

# Expose no ports (this is a background service)

# Set entrypoint
ENTRYPOINT ["./cloudflare-backuper"]

# Default command (can be overridden with -once for single backup)
CMD ["-config", "/app/config.yml"]
