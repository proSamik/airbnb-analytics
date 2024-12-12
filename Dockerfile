# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod and sum files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install necessary runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    curl \
    postgresql-client \
    && update-ca-certificates

# Copy only necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/go.mod /app/go.sum ./

# Make init script executable
COPY --from=builder /app/scripts/startup/init.sh .
RUN chmod +x init.sh

# Set environment variables if needed
ENV APP_ENV=production

# Expose port if your application uses a specific port
# EXPOSE 8080

# Use shell form to allow init script to handle multiple commands
ENTRYPOINT ["./init.sh"]