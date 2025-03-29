# Simple REST API

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

#### User Endpoints
- `GET /api/v1/users/{id}` - Get a user by ID
- `POST /api/v1/users` - Create a new user

#### Task Endpoints (Employer Role)
- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks` - Get all tasks with filtering and sorting
- `GET /api/v1/employee-summary` - Get task statistics for all employees

#### Task Endpoints (Employee Role)
- `GET /api/v1/tasks/assigned` - Get tasks assigned to the employee
- `PATCH /api/v1/tasks/{id}/status` - Update task status

For detailed API documentation with request/response examples, see:
- [User Service API Documentation](docs/user_service.md)
- [Task Service API Documentation](docs/task_service.md)

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (for containerization)
- Make (optional, for running Makefile commands)

### Quick Start (with Make)

The easiest way to get started is using the `setup` command provided in the Makefile:

```bash
# Clone the repository
git clone https://github.com/llkhacquan/cisab.git
cd cisab

# Run the setup command
make setup
```

This will:
1. Stop any existing containers
2. Start the PostgreSQL database
3. Run database migrations

### Running the Tests

Run all tests and linting:

```bash
make check
```

This command will:
1. Format code using goimports
2. Run linters (golangci-lint, pre-commit hooks)
3. Run all tests

### Using Docker

#### Running with Docker Compose (recommended)

```bash
# Build and start all services (PostgreSQL, API)
docker-compose up -d

# To stop all services
docker-compose down
```

The API will be available at http://localhost:8080

#### Running Only the Database

If you want to run the PostgreSQL database in Docker but develop the API locally:

```bash
# Start only the database
make db-up

# Run migrations
go run cmd/migrate/main.go

# Run the API locally
go run main.go
```

#### Building and Running the API Container Manually

```bash
# Build the Docker image
make build

# Run the container
make run
```

### Manual Setup (without Docker)

If you prefer to run everything locally:

1. Install PostgreSQL locally
2. Create a database named `cisab`
3. Update the database configuration in `config/local.yaml` to point to your local database
4. Run migrations:
   ```bash
   go run cmd/migrate/main.go
   ```
5. Start the API:
   ```bash
   go run main.go
   ```

### Database Migrations

Before running the application for the first time, you need to run the database migrations:

```bash
go run cmd/migrate/main.go
```

This will create all the necessary tables in the database.

### Project Structure

```
cisab/
├── cmd/
│   └── migrate/        # Database migration tool
├── config/             # Configuration files
├── db/                 # Database migrations
├── docs/               # API documentation
├── pkg/
│   ├── api/            # HTTP API handlers and middleware
│   ├── authctx/        # Authentication context
│   ├── config/         # Configuration loading
│   ├── dbctx/          # Database context
│   ├── models/         # Data models
│   ├── repo/           # Data access layer
│   ├── service/        # Business logic
│   ├── testutil/       # Testing utilities
│   └── utils/          # Utility functions
├── scripts/            # Helper scripts
├── docker-compose.yaml # Docker Compose configuration
├── Dockerfile          # Docker build configuration
├── Makefile            # Make commands
└── README.md           # This file
```

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

### Update task status (protected endpoint)

```bash
curl -X PATCH http://localhost:8080/api/v1/tasks/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
    "status": "completed"
  }'
```

### Get tasks assigned to an employee (protected endpoint)

```bash
curl -X GET "http://localhost:8080/api/v1/tasks/assigned?status=pending&sort_by=created_at&sort_order=desc" \
  -H "Authorization: Bearer your-jwt-token"
```

### Get all tasks for an employer (protected endpoint)

```bash
curl -X GET "http://localhost:8080/api/v1/tasks?status=pending&assignee_id=3&sort_by=due_date" \
  -H "Authorization: Bearer your-jwt-token"
```

### Get employee task summary (protected endpoint)

```bash
curl -X GET http://localhost:8080/api/v1/employee-summary \
  -H "Authorization: Bearer your-jwt-token"
```

For more examples and detailed API documentation, see:
- [User Service API Documentation](docs/user_service.md)
- [Task Service API Documentation](docs/task_service.md)

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
  name: cisab
```

### Environment Variables

Configuration values can be overridden using environment variables:

- `PORT` - The port the server will listen on (overrides `server.port` in the config file)
- `ENVIRONMENT` - The current environment (e.g., `local`, `dev`, `staging`, `production`)
- `DATABASE_HOST` - Database hostname (overrides `database.host` in the config file)
- `DATABASE_PORT` - Database port (overrides `database.port` in the config file)
- `DATABASE_USER` - Database username (overrides `database.user` in the config file)
- `DATABASE_PASSWORD` - Database password (overrides `database.password` in the config file)
- `DATABASE_NAME` - Database name (overrides `database.name` in the config file)
- `JWT_SECRET` - Secret key for JWT token signing (overrides `jwt.secret` in the config file)
- `JWT_TTL_SECOND` - JWT token expiration time in seconds (overrides `jwt.ttl_second` in the config file)

You can also create a `.env` file in the project root to set these variables.

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you're having trouble connecting to the database:

1. Verify PostgreSQL is running:
   ```bash
   docker ps
   ```

2. Check the database logs:
   ```bash
   docker-compose logs db
   ```

3. Make sure the database configuration matches your environment:
   - When running locally: use `localhost` as the host and `5433` as the port
   - When running in Docker: use `db` as the host and `5432` as the port (within the Docker network)

#### Migration Issues

If database migrations fail:

1. Make sure the database exists and is accessible
2. Check the migration logs for specific errors:
   ```bash
   go run cmd/migrate/main.go -v
   ```

#### API Server Issues

If the API server fails to start:

1. Check the logs for error messages
2. Verify that all required environment variables are set
3. Make sure the database is running and migrations have been applied
4. Check that the port (default: 8080) is not already in use

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
