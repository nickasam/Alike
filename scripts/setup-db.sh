#!/bin/bash

set -e

echo "🚀 Setting up Alike Database..."

# Check PostgreSQL installation
if ! command -v psql &> /dev/null; then
    echo "⚠️  PostgreSQL not found in PATH"
    echo "📦 Installing PostgreSQL..."
    brew install postgresql@15
    
    echo "⏳ Waiting for installation to complete..."
    sleep 10
    
    # Try to find PostgreSQL
    if [ -f /opt/homebrew/bin/psql ]; then
        export PATH="/opt/homebrew/bin:$PATH"
    fi
fi

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL installation failed or not in PATH"
    echo "Please try: brew install postgresql@15"
    echo "Then add to PATH: export PATH=\"/opt/homebrew/bin:\$PATH\""
    exit 1
fi

echo "✅ PostgreSQL found: $(psql --version)"

# Start PostgreSQL service
echo "🔧 Starting PostgreSQL service..."
brew services start postgresql@15 2>/dev/null || brew services restart postgresql@15

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to start..."
sleep 5

# Initialize database if needed
if ! pg_isready -q; then
    echo "⚠️  PostgreSQL not ready, initializing..."
    initdb -D /opt/homebrew/var/postgresql@15 2>/dev/null || true
    sleep 3
fi

# Create database
echo "📊 Creating database..."
createdb alike_db 2>/dev/null && echo "✅ Database created" || echo "ℹ️  Database already exists"

# Create user
echo "👤 Creating database user..."
psql -d postgres << 'EOSQL' 2>/dev/null || echo "User already exists"
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'alike_user') THEN
        CREATE USER alike_user WITH PASSWORD 'alike_password';
    END IF;
    GRANT ALL PRIVILEGES ON DATABASE alike_db TO alike_user;
    ALTER USER alike_user WITH PASSWORD 'alike_password';
END
$$;
EOSQL

# Test connection
echo "🔗 Testing connection..."
if PGPASSWORD=alike_password psql -U alike_user -d alike_db -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Database connection successful!"
else
    echo "⚠️  Connection test failed - PostgreSQL may still be starting"
    echo "   Try again in 10 seconds or check: brew services list"
fi

# Run migrations
echo "📋 Running database migrations..."
cd /Users/zhenghongfei6/go/src/github.com/Alike
if go run cmd/migrate/main.go up 2>&1; then
    echo "✅ Migrations completed"
else
    echo "⚠️  Migrations failed - you may need to run manually"
fi

# Seed test data
echo "🌱 Seeding test data..."
if PGPASSWORD=alike_password psql -U alike_user -d alike_db -f db/seeds/seed.sql > /dev/null 2>&1; then
    echo "✅ Test data seeded"
else
    echo "ℹ️  Seed data already exists or failed"
fi

echo ""
echo "✅ Database setup complete!"
echo ""
echo "📊 Connection Info:"
echo "   Host: localhost"
echo "   Port: 5432"
echo "   Database: alike_db"
echo "   User: alike_user"
echo "   Password: alike_password"
echo ""
echo "🧪 Test connection:"
echo "   PGPASSWORD=alike_password psql -U alike_user -d alike_db"
echo ""
echo "🚀 Next:"
echo "   go run cmd/api/main.go"
echo "   open web/public/index.html"
echo ""
echo "🎉 Ready to use!"
