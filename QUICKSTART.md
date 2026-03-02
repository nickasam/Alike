# Alike - Quick Start Guide

## Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Redis 7+ (optional, for caching)
- Docker & Docker Compose (optional, for easy deployment)

---

## Quick Start with Docker (Recommended)

### 1. Clone the Repository
```bash
git clone https://github.com/nickasam/Alike.git
cd Alike
```

### 2. Configure Environment
```bash
# Copy example configuration
cp config/config.yaml.example config/config.yaml

# Edit configuration (optional for development)
vim config/config.yaml
```

### 3. Start Services
```bash
# Start PostgreSQL, Redis, and API
docker-compose -f deployments/docker/docker-compose.yml up -d

# Or use the deployment script
./scripts/deploy.sh
```

### 4. Run Migrations
```bash
# Run database migrations
docker-compose -f deployments/docker/docker-compose.yml exec api \
  go run cmd/migrate/main.go up
```

### 5. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"ok"}
```

### 6. View Logs
```bash
docker-compose -f deployments/docker/docker-compose.yml logs -f api
```

---

## Manual Setup (Without Docker)

### 1. Install Dependencies
```bash
# Install Go 1.23+
# Install PostgreSQL 15+
# Install Redis 7+ (optional)

# Install Go dependencies
go mod download
```

### 2. Set Up Database
```bash
# Create database
createdb alike_db

# Create user
psql -d postgres -c "CREATE USER alike_user WITH PASSWORD 'your_password';"
psql -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE alike_db TO alike_user;"
```

### 3. Configure Environment
```bash
# Copy example configuration
cp config/config.yaml.example config/config.yaml

# Edit with your database credentials
vim config/config.yaml
```

### 4. Run Migrations
```bash
go run cmd/migrate/main.go up
```

### 5. Build and Run
```bash
# Build
make build

# Or run directly
go run cmd/api/main.go

# Or with the Makefile
make run
```

### 6. Test the API
```bash
curl http://localhost:8080/health
```

---

## Development

### Start Development Server with Hot Reload
```bash
# Using the dev script
./scripts/dev.sh

# Or manually with air
air

# Or without hot reload
go run cmd/api/main.go
```

### Run Tests
```bash
make test
```

### Build for Production
```bash
make build
# Binary will be in ./bin/alike-api
```

---

## API Testing

### Using cURL

#### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+8613800138000",
    "verification_code": "123456",
    "nickname": "Test User",
    "password": "password123"
  }'
```

#### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "+8613800138000",
    "password": "password123"
  }'
```

Save the `access_token` from the response for next requests.

#### 3. Get Nearby Users
```bash
curl -X GET "http://localhost:8080/api/v1/users/nearby?lat=31.2304&lng=121.4737&radius=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Using Postman

1. Import the API collection (if available)
2. Set environment variable `base_url` to `http://localhost:8080`
3. Run authentication requests first
4. Use the received token for authenticated requests

---

## Common Issues

### Database Connection Failed
- Check PostgreSQL is running: `pg_isready`
- Verify database credentials in `config/config.yaml`
- Ensure database exists: `psql -l | grep alike_db`

### Port Already in Use
- Change port in `config/config.yaml`
- Or kill process using port 8080: `lsof -ti:8080 | xargs kill -9`

### Migration Errors
- Drop and recreate database: `dropdb alike_db && createdb alike_db`
- Run migrations again: `go run cmd/migrate/main.go up`

---

## Production Deployment

### Security Checklist
- [ ] Change JWT secret
- [ ] Use strong database password
- [ ] Enable HTTPS
- [ ] Set up firewall
- [ ] Enable rate limiting
- [ ] Configure backups
- [ ] Set up monitoring

### Deploy with Docker
```bash
# Build production image
docker build -f deployments/docker/Dockerfile -t alike-api:latest .

# Run with docker-compose
docker-compose -f deployments/docker/docker-compose.yml up -d
```

### Deploy to Server
1. Copy files to server
2. Install dependencies
3. Configure environment
4. Run migrations
5. Start service with systemd or PM2

---

## Next Steps

1. **Read API Documentation**: Check `API.md` for all endpoints
2. **Set Up Mobile App**: Start iOS/Android development
3. **Configure Push Notifications**: Set up FCM/APNs
4. **Set Up Monitoring**: Configure Prometheus + Grafana
5. **Configure Domain**: Set up custom domain and SSL

---

## Support

- GitHub Issues: https://github.com/nickasam/Alike/issues
- API Docs: See `API.md`
- Architecture: See `docs/architecture.md`

---

**Happy Coding! 🚀**
