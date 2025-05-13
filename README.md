# E-Commerce Platform (Microservices)

## Overview

This project is a basic e-commerce platform built with Go, structured using Clean Architecture and composed of three microservices:

- API Gateway – routes incoming requests to the appropriate service, handles logging and basic authentication.

- Inventory Service – manages product data, categories, stock, and prices.

- Order Service – handles order creation, status updates, and product quantities per order.

All services are written in Go using the Gin framework and connected to PostgreSQL for persistence.

## Project Structure

```
.
├── inventory/              # Inventory microservice
│   ├── cmd/
│   └── internal/
├── order/                 # Order microservice
│   ├── cmd/
│   └── internal/
└── gateway/               # API Gateway
├── cmd/
└── internal/
```

## New Files Added for Redis Implementation

During the Redis caching implementation, several new files were added to enhance the system's performance:

### API Gateway Changes

1. `apiGateway/internal/cache/redis.go`
   - Implements Redis client connection and configuration
   - Handles Redis connection pooling and error handling
   - Provides methods for cache operations (Get, Set, Delete)

2. `apiGateway/internal/cache/user_cache.go`
   - Implements user-specific caching logic
   - Defines cache key formats and TTL settings
   - Provides methods for user profile caching operations

3. `apiGateway/internal/handler/user_handler.go` (modified)
   - Updated to integrate Redis caching
   - Implements cache-first strategy for user profile retrieval
   - Handles cache invalidation on profile updates

### Configuration Changes

1. `apiGateway/config/config.go` (modified)
   - Added Redis configuration settings
   - Includes Redis connection parameters (host, port, password)
   - Defines cache-related constants

### Dependencies

1. `apiGateway/go.mod` (modified)
   - Added Redis client dependency: `github.com/redis/go-redis/v9`
   - Updated other dependencies to support Redis integration

These new files work together to provide a robust caching layer that:
- Reduces database load
- Improves response times
- Maintains data consistency
- Handles cache invalidation
- Provides graceful fallback to database when cache is unavailable

## Requirements

- Go 1.18 or higher

- PostgreSQL

- Redis (for caching)

- Basic understanding of RESTful APIs

## Redis Caching

The project uses Redis for caching frequently accessed data, particularly user profiles. This improves performance by reducing database load and response times for read-heavy operations.

### Redis Configuration
- Default port: 6379
- Cache TTL: 30 minutes
- Keys format: `user:{id}`

### Testing Redis Caching

1. Start Redis server:
```bash
redis-server
```

2. Test caching flow using Postman or curl:

a. First profile request (cache miss):
```bash
curl -X GET http://localhost:8080/profile/5 \
  -H "Authorization: Basic YWRtaW46MTIzNA=="
```

b. Second profile request (cache hit):
```bash
curl -X GET http://localhost:8080/profile/5 \
  -H "Authorization: Basic YWRtaW46MTIzNA=="
```

c. Update profile (invalidates cache):
```bash
curl -X PUT http://localhost:8080/profile/5 \
  -H "Authorization: Basic YWRtaW46MTIzNA==" \
  -H "Content-Type: application/json" \
  -d '{"username": "updateduser"}'
```

d. Third profile request (cache miss):
```bash
curl -X GET http://localhost:8080/profile/5 \
  -H "Authorization: Basic YWRtaW46MTIzNA=="
```

### Redis CLI Commands

Monitor cache operations:
```bash
# Connect to Redis CLI
redis-cli

# List all keys
KEYS *

# Get specific user data
GET user:5

# Check TTL of a key
TTL user:5

# Monitor Redis operations in real-time
MONITOR
```

### Testing Redis with Postman

1. **Setup Postman Collection**
   - Create a new collection named "E-Commerce API"
   - Add the following environment variables:
     - `base_url`: http://localhost:8080
     - `auth`: Basic YWRtaW46MTIzNA==

2. **Test Cache Miss (First Request)**
   - Create a new GET request
   - URL: `{{base_url}}/profile/5`
   - Headers:
     - Authorization: {{auth}}
   - Expected: Slower response time (data from database)

3. **Test Cache Hit (Second Request)**
   - Use the same GET request
   - URL: `{{base_url}}/profile/5`
   - Headers:
     - Authorization: {{auth}}
   - Expected: Faster response time (data from Redis)

4. **Test Cache Invalidation**
   - Create a new PUT request
   - URL: `{{base_url}}/profile/5`
   - Headers:
     - Authorization: {{auth}}
     - Content-Type: application/json
   - Body (raw JSON):
     ```json
     {
         "username": "updateduser"
     }
     ```
   - Expected: Cache is invalidated

5. **Verify Cache Miss After Update**
   - Use the GET request again
   - URL: `{{base_url}}/profile/5`
   - Headers:
     - Authorization: {{auth}}
   - Expected: Slower response time (data from database)

6. **Verify Cache Hit After Update**
   - Use the GET request one more time
   - URL: `{{base_url}}/profile/5`
   - Headers:
     - Authorization: {{auth}}
   - Expected: Faster response time (data from Redis)

7. **Monitor Redis in Real-time**
   - Open Redis CLI in a separate terminal
   - Run `MONITOR` command
   - Execute the Postman requests
   - Observe Redis operations in the terminal

Expected Results:
- First GET request: ~100-200ms (database)
- Second GET request: ~10-20ms (Redis)
- After PUT request: ~100-200ms (database)
- Final GET request: ~10-20ms (Redis)

## Installation

### 1. Install Go Dependencies:

   In each service directory:
   ```
    cd inventoryService
    go mod tidy

    cd ../orderService
    go mod tidy

    cd ../apiGateway
    go mod tidy
   ```

### 2. Set Up PostgreSQL Tables:
Inventory Service:
```
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL NOT NULL,
    stock INT NOT NULL,
    category_id INT NOT NULL
);
```

Order Service:
```
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL,
    quantity INT NOT NULL
);
```

### 3. Running the services:
Inventory Service:
```
cd inventoryService
go run cmd/main.go
```
Runs on http://localhost:8081

Order Service:
```
cd orderService
go run cmd/main.go
```
Runs on http://localhost:8082

API-Gateway:
```
cd apiGateway
go run cmd/main.go
```
Runs on http://localhost:8080

## Authentication
All API Gateway endpoints are protected by Basic Auth:

- Username: admin

- Password: 1234

Base64 encoded: YWRtaW46MTIzNA==

Include this in your Authorization header:
```
Authorization: Basic YWRtaW46MTIzNA==
```

## Sample API Requests (via API Gateway)

### Create Product:
```
curl -X POST http://localhost:8080/products \
  -H "Authorization: Basic YWRtaW46MTIzNA==" \
  -H "Content-Type: application/json" \
  -d '{"name": "iPhone", "description": "Apple smartphone", "price": 999.99, "stock": 10, "category_id": 1}'
```

### List Products:
```
curl -X GET http://localhost:8080/products \
  -H "Authorization: Basic YWRtaW46MTIzNA=="
```

### Create Order:
```
curl -X POST http://localhost:8080/orders \
  -H "Authorization: Basic YWRtaW46MTIzNA==" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "items": [
      {"product_id": 1, "quantity": 2}
    ]
  }'
```

### Get Order by ID:
```
curl -X GET http://localhost:8080/orders/1 \
  -H "Authorization: Basic YWRtaW46MTIzNA=="
```

### Update Order Status:
```
curl -X PATCH http://localhost:8080/orders/1 \
  -H "Authorization: Basic YWRtaW46MTIzNA==" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
```

## Logging
Each request is logged to the console with:

- IP address

- Method

- Endpoint

- Status code

- Duration


## HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Create successful
- `400 Bad Request`: Missing required parameters or invalid input
- `401 Unauthorized (Auth)`: Unauthorized

- `404 Not Found`: No news found for the specified cryptocurrency
- `500 Internal Server Error`: Server-side error

## Notes
- All API calls must go through the API Gateway (localhost:8080)

- Services run independently; no Docker or gRPC is used

- You can easily extend this to include user authentication, payment providers, etc.
