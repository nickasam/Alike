#!/bin/bash

set -e

echo "🚀 Deploying Alike..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Create config file if it doesn't exist
if [ ! -f config/config.yaml ]; then
    echo "📝 Creating config file..."
    mkdir -p config
    cp config/config.yaml.example config/config.yaml
    echo "⚠️  Please edit config/config.yaml with your settings before deploying!"
    exit 0
fi

# Build and start services
echo "🏗️  Building Docker images..."
docker-compose -f deployments/docker/docker-compose.yml build

echo "🚀 Starting services..."
docker-compose -f deployments/docker/docker-compose.yml up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Run migrations
echo "📊 Running database migrations..."
docker-compose -f deployments/docker/docker-compose.yml exec -T api go run cmd/migrate/main.go up

echo "✅ Deployment complete!"
echo ""
echo "🌐 API is available at: http://localhost:8080"
echo "❤️  Health check: http://localhost:8080/health"
echo ""
echo "📝 View logs with: docker-compose -f deployments/docker/docker-compose.yml logs -f"
echo "🛑 Stop services with: docker-compose -f deployments/docker/docker-compose.yml down"
