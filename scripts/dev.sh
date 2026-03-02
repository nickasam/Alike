#!/bin/bash

# Development script with hot reload

echo "🔧 Starting development server..."

# Install air if not installed
if ! command -v air &> /dev/null; then
    echo "📦 Installing air for hot reload..."
    go install github.com/cosmtrek/air@latest
fi

# Run database migrations
echo "📊 Running migrations..."
go run cmd/migrate/main.go up

# Start development server with hot reload
echo "🚀 Starting server..."
air || go run cmd/api/main.go
