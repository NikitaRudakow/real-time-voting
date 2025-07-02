# Real-Time Voting System

## 🏗️ Architecture

```
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers, routing, middleware
│   ├── config/         # Configuration management
│   ├── database/       # DB connection and migrations
│   ├── logger/         # Logging configuration
│   ├── models/         # Data models
│   ├── dto/            # Request/response DTOs
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── websocket/      # Real-time communication
├── migrations/         # Database migrations
├── docs/               # Swagger documentation
└── docker-compose.yml  # Container orchestration
```

## 🛠️ Technology Stack

- **Backend**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15
- **Authentication**: JWT with bcrypt password hashing
- **Real-time**: WebSocket (Gorilla WebSocket)
- **Documentation**: Swagger/OpenAPI
- **Logging**: Zerolog (middleware-based)
- **Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (for containerized deployment)
- Make (optional, for using Makefile commands)

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd real-time-voting
   ```

2. **Start the application**
   ```bash
   docker-compose up --build
   ```

3. **Access the API**
   - API Base URL: `http://localhost:8080`
   - Health Check: `http://localhost:8080/health`
   - Swagger Documentation: `http://localhost:8080/swagger/index.html`
   - WebSocket: `ws://localhost:8080/api/v1/ws`

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   docker run -d --name postgres \
     -e POSTGRES_DB=voting_db \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 \
     postgres:15-alpine
   ```

3. **Run migrations**
   ```bash
   make migrate-up
   ```

4. **Generate Swagger docs**
   ```bash
   swag init -g cmd/server/main.go
   ```

5. **Start the application**
   ```bash
   make run
   ```

## 🔐 Authentication

The system uses JWT (JSON Web Tokens) for authentication. Example usage:

### 1. Register a new user
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login to get JWT token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }'
```

### 3. Use the token for authenticated requests
```bash
curl -X POST http://localhost:8080/api/v1/polls \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "question": "Pineapple on pizza: Yes or No?",
    "options": ["Yes", "No", "Maybe"]
  }'
```

## 📚 API Documentation

- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Base URL**: `http://localhost:8080/api/v1`

### Endpoints

#### Authentication (Public)
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login user

#### Health Check
- `GET /health` - Server health status

#### Polls
- `GET /polls` - Get all polls (Public)
- `GET /polls/active` - Get active polls (Public)
- `GET /polls/{id}` - Get poll by ID (Public)
- `POST /polls` - Create poll (Authenticated)
- `PUT /polls/{id}` - Update poll (Authenticated)
- `DELETE /polls/{id}` - Delete poll (Authenticated)
- `POST /polls/{id}/vote` - Cast vote (Authenticated)

#### WebSocket
- `GET /ws` - WebSocket connection for real-time updates

### Authentication Headers
For authenticated endpoints, include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## 🔧 Configuration

The application can be configured using environment variables:

| Variable        | Default                                                        | Description                |
|----------------|----------------------------------------------------------------|----------------------------|
| `PORT`         | `8080`                                                         | Server port                |
| `DATABASE_URL` | `postgres://postgres:password@localhost:5432/voting_db?sslmode=disable` | Database connection string |
| `LOG_LEVEL`    | `info`                                                         | Logging level              |
| `JWT_SECRET`   | `your-secret-key`                                              | Secret key for JWT tokens  |

## 🧩 Middleware

- **AuthMiddleware**: Verifies JWT and adds user_id to the request context.
- **OptionalAuth**: Adds user_id to the context if a token is present, but does not require it.
- **LoggingMiddleware**: Centrally logs all HTTP requests (method, path, status, user_id, query, request body, execution time).
- **ErrorHandlingMiddleware**: Catches and logs all errors and panics, returns a standard error response.

## 📦 Development Commands

```bash
# Build the application
make build

# Run locally
make run


# Format code
make fmt

# Lint code
make lint

# Clean build artifacts
make clean

# Install dependencies
make deps

# Update go.mod
go mod tidy

# Generate Swagger docs
swag init -g cmd/server/main.go
```

## 🐳 Docker Commands

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-run

# Stop Docker Compose
make docker-stop
```

## 📊 Database Schema

### Tables

**users**
- `id` (UUID, Primary Key)
- `username` (VARCHAR(50), UNIQUE, NOT NULL)
- `email` (VARCHAR(255), UNIQUE, NOT NULL)
- `password` (VARCHAR(255), NOT NULL) - bcrypt hashed
- `created_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)
- `updated_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)

**polls**
- `id` (UUID, Primary Key)
- `question` (TEXT, NOT NULL)
- `is_active` (BOOLEAN, NOT NULL, DEFAULT true)
- `created_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)
- `updated_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)
- `ends_at` (TIMESTAMP WITH TIME ZONE, NULLABLE)
- `created_by` (UUID, Foreign Key to users.id)

**poll_options**
- `id` (UUID, Primary Key)
- `poll_id` (UUID, Foreign Key to polls.id)
- `text` (TEXT, NOT NULL)
- `created_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)

**votes**
- `id` (UUID, Primary Key)
- `poll_id` (UUID, Foreign Key to polls.id)
- `option_id` (UUID, Foreign Key to poll_options.id)
- `voter_id` (UUID, Foreign Key to users.id)
- `created_at` (TIMESTAMP WITH TIME ZONE, NOT NULL)