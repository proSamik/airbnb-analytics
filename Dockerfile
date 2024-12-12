# Build stage
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache go

# Copy the binary and scripts
COPY --from=builder /app/main .
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/go.* ./

# Make init script executable
COPY scripts/startup/init.sh .
RUN chmod +x init.sh

# Command to run init script
CMD ["./init.sh"]