# Knovel Simple REST API

A simple HTTP REST API with middleware demonstration.

## Features

- HTTP REST API with Gorilla Mux
- Multiple middleware layers:
  - Logging middleware (logs all requests with timing)
  - CORS middleware (handles CORS headers)
  - Authentication middleware (mock the authentication process now)
- Structured JSON responses
- Environment variable configuration

## API Endpoints

### Public Endpoints

- `GET /health` - Health check endpoint

### Protected Endpoints (require X-API-Key header)

- `GET /api/v1/users` - Get a list of users

## Getting Started

### Prerequisites

- Go 1.24 or later

### Running the API

1. Clone the repository
2. Install dependencies:
   ```
   go mod tidy
   ```
3. Run the server:
   ```
   go run main.go
   ```

The server will start on port 8080 by default. You can configure the port by setting the `PORT` environment variable.

## Using the API

### Health Check

```
curl -X GET http://localhost:8080/health
```

### Get Users (protected endpoint)

```
curl -X GET http://localhost:8080/api/v1/users -H "X-API-Key: your-api-key"
```

## Configuration

The application uses a configuration system based on YAML files. Configuration files are located in the `config/` directory and are named after the environment (e.g., `local.yaml`, `production.yaml`).

### Configuration Files

The default configuration file is `config/local.yaml`. You can create different configuration files for different environments.

Example configuration file:

```yaml
# Local environment configuration
environment: local

# Server configuration
server:
  port: 8080
```

### Environment Variables

Configuration values can be overridden using environment variables:

- `PORT` - The port the server will listen on (overrides `server.port` in the config file)
- `ENVIRONMENT` - The current environment (e.g., `local`, `dev`, `staging`, `production`)

You can also create a `.env` file in the project root to set these variables.
