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

For detailed API documentation with request/response examples,
see [User Service API Documentation](docs/user_service.md).

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
