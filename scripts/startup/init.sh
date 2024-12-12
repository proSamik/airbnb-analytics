#!/bin/sh

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Run database setup
echo "Running database setup..."
go run scripts/db_setup.go

# Start the main application
echo "Starting main application..."
./main
