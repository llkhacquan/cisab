# Knovel Simple REST API

A simple HTTP REST API with middleware demonstration.

## Features

- HTTP REST API with Gorilla Mux
- Multiple middleware layers:
    - Logging middleware (logs all requests with timing)
    - CORS middleware (handles CORS headers)
    - Authentication middleware (JWT-based authentication)
- User management (registration, authentication)
- Structured JSON responses
- Environment variable configuration
- PostgreSQL database integration with GORM

## API Endpoints

### Public Endpoints

- `GET /health` - Health check endpoint
- `POST /api/v1/login` - Authenticate and get JWT token

### Protected Endpoints (require JWT authentication)

- `GET /api/v1/users/{id}` - Get a user by ID
- `POST /api/v1/users` - Create a new user
- `POST /api/v1/tasks` - Create a new task

For detailed API documentation with request/response examples,
see [API Documentation](docs/user_service.md).

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (for running PostgreSQL)

### Running the test

```bash
make check
```

### Setting up the database

The project uses PostgreSQL as the database. You can easily set it up using Docker Compose:

```bash
docker-compose up -d
```

This will start PostgreSQL on port 5433.

### Database migrations

Before running the application for the first time, you need to run the database migrations:

```bash
go run cmd/migrate/main.go
```

### Running the API

1. Clone the repository
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Start the database:
   ```
   docker-compose up -d
   ```
4. Run the server:
   ```
   go run main.go
   ```

The server will start on port 8080 by default. You can configure the port by setting the `PORT` environment variable.

## Using the API

### Health Check

```bash
curl -X GET http://localhost:8080/health
```

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword",
    "role": "employee"
  }'
```

### Login and get JWT token

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword"
  }'
```

### Get user by ID (protected endpoint)

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer your-jwt-token"
```

### Create a new task (protected endpoint)

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API",
    "due_date": "2023-04-15T00:00:00Z"
  }'
```

For more examples and detailed API documentation, see [API Documentation](docs/user_service.md).

## Configuration

The application uses a configuration system based on YAML files. Configuration files are located in the `config/`
directory and are named after the environment (e.g., `local.yaml`, `production.yaml`).

### Configuration Files

The default configuration file is `config/local.yaml`. You can create different configuration files for different
environments.

Example configuration file:

```yaml
# Local environment configuration
environment: local

# Server configuration
server:
  port: 8080

# Database configuration
database:
  host: localhost
  port: 5433
  user: postgres
  password: password
  name: knovel
```

### Environment Variables

Configuration values can be overridden using environment variables:

- `PORT` - The port the server will listen on (overrides `server.port` in the config file)
- `ENVIRONMENT` - The current environment (e.g., `local`, `dev`, `staging`, `production`)

You can also create a `.env` file in the project root to set these variables.
