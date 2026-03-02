# Alike API Documentation

## Base URL
```
Production: https://api.alike.app
Development: http://localhost:8080
```

## Authentication

Most endpoints require authentication via JWT token.

### Headers
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

## Endpoints

### Health Check
```http
GET /health
```

**Response (200):**
```json
{
  "status": "ok"
}
```

---

### Authentication

#### Register
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "phone": "+8613800138000",
  "verification_code": "123456",
  "nickname": "John",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "phone": "+8613800138000",
      "nickname": "John",
      "created_at": "2026-03-02T15:00:00Z"
    },
    "tokens": {
      "access_token": "eyJhbG...",
      "refresh_token": "eyJhbG..."
    }
  }
}
```

#### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "phone": "+8613800138000",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "user": {...},
    "tokens": {...}
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "eyJhbG..."
}
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG..."
  }
}
```

#### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

---

### Users

#### Get Current User
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "phone": "+8613800138000",
    "nickname": "John",
    "avatar_url": "https://...",
    "bio": "Hello!",
    "location_lat": 31.2304,
    "location_lng": 121.4737
  }
}
```

#### Update Current User
```http
PUT /api/v1/users/me
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "nickname": "New Nickname",
  "bio": "Updated bio",
  "location_lat": 31.2304,
  "location_lng": 121.4737
}
```

#### Get Nearby Users
```http
GET /api/v1/users/nearby?lat=31.2304&lng=121.4737&radius=10&page=1&limit=20
Authorization: Bearer <token>
```

**Query Parameters:**
- `lat` (float, required): Latitude
- `lng` (float, required): Longitude
- `radius` (float, optional): Radius in km (default: 10)
- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 20)

**Response (200):**
```json
{
  "success": true,
  "data": {
    "users": [...],
    "page": 1,
    "limit": 20
  }
}
```

---

### Matches

#### List Matches
```http
GET /api/v1/matches
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "user1_id": "uuid",
      "user2_id": "uuid",
      "last_message_at": "2026-03-02T15:00:00Z"
    }
  ]
}
```

#### Get Match
```http
GET /api/v1/matches/:id
Authorization: Bearer <token>
```

#### Like User
```http
POST /api/v1/matches/:id/like
Authorization: Bearer <token>
```

---

### Chats

#### List Chats
```http
GET /api/v1/chats
Authorization: Bearer <token>
```

#### Get Chat
```http
GET /api/v1/chats/:id
Authorization: Bearer <token>
```

#### Get Messages
```http
GET /api/v1/chats/:id/messages?page=1&limit=50
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "chat_id": "uuid",
      "sender_id": "uuid",
      "content": "Hello!",
      "message_type": "text",
      "is_read": false,
      "created_at": "2026-03-02T15:00:00Z"
    }
  ]
}
```

#### Send Message
```http
POST /api/v1/chats/:id/messages
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "content": "Hello!",
  "type": "text"
}
```

---

### Notifications

#### List Notifications
```http
GET /api/v1/notifications?page=1&limit=20
Authorization: Bearer <token>
```

#### Mark as Read
```http
POST /api/v1/notifications/:id/read
Authorization: Bearer <token>
```

#### Mark All as Read
```http
POST /api/v1/notifications/read-all
Authorization: Bearer <token>
```

---

## Error Responses

All errors follow this format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` - Invalid request data
- `UNAUTHORIZED` - Missing or invalid token
- `NOT_FOUND` - Resource not found
- `INTERNAL_ERROR` - Server error
- `USER_NOT_FOUND` - User does not exist
- `INVALID_PASSWORD` - Wrong password

### HTTP Status Codes

- `200` - Success
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

---

## Testing with cURL

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+8613800138000","password":"password123"}'
```

### Get Nearby Users
```bash
curl -X GET "http://localhost:8080/api/v1/users/nearby?lat=31.2304&lng=121.4737&radius=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Send Message
```bash
curl -X POST http://localhost:8080/api/v1/chats/CHAT_ID/messages \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Hello!","type":"text"}'
```

---

## Rate Limiting

- 60 requests per minute per IP
- Burst: 10 requests

Exceeding limits returns `429 Too Many Requests`.

---

## WebSocket

WebSocket connection for real-time messaging:

```
ws://localhost:8080/api/v1/chats/:id/ws?token=<jwt_token>
```

Messages are sent/received as JSON:

```json
{
  "type": "message",
  "data": {
    "id": "uuid",
    "content": "Hello!",
    "sender_id": "uuid",
    "created_at": "2026-03-02T15:00:00Z"
  }
}
```

---

*Last updated: 2026-03-02*
*Version: 1.0.0*
