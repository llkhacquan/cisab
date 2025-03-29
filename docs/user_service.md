# User Service API Documentation

This document provides detailed information about the User Service API endpoints, including request/response formats and
curl examples.

## Table of Contents

1. [User Registration](#user-registration)
2. [User Login](#user-login)
3. [Get User by ID](#get-user-by-id)
4. [Get Current User](#get-current-user)

## User Registration

Register a new user in the system.

### Endpoint

```
POST /api/v1/users
```

### Request Body

```json
{
  "name": "string",
  "email": "string",
  "password": "string",
  "role": "employee|employer"
}
```

#### Fields

| Field    | Type   | Required | Description                                    |
|----------|--------|----------|------------------------------------------------|
| name     | string | Yes      | User's full name                               |
| email    | string | Yes      | User's email address (must be unique)          |
| password | string | Yes      | User's password (minimum 8 characters)         |
| role     | string | Yes      | User's role - must be "employee" or "employer" |

### Response

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "employee",
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
}
```

### Example

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

### Error Responses

| Status Code | Error Message           | Description                                |
|-------------|-------------------------|--------------------------------------------|
| 400         | Missing required fields | One or more required fields are missing    |
| 400         | Invalid role            | Role must be "employee" or "employer"      |
| 400         | Password too short      | Password must be at least 8 characters     |
| 409         | User already exists     | A user with this email already exists      |
| 500         | Internal server error   | An unexpected error occurred on the server |

## User Login

Authenticate a user and get a JWT token.

### Endpoint

```
POST /api/v1/login
```

### Request Body

```json
{
  "email": "string",
  "password": "string"
}
```

#### Fields

| Field    | Type   | Required | Description          |
|----------|--------|----------|----------------------|
| email    | string | Yes      | User's email address |
| password | string | Yes      | User's password      |

### Response

```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "employee",
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    },
    "token_expiry": 1680355200
  }
}
```

### Example

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword"
  }'
```

### Error Responses

| Status Code | Error Message           | Description                                |
|-------------|-------------------------|--------------------------------------------|
| 400         | Missing required fields | One or more required fields are missing    |
| 401         | Invalid credentials     | Email or password is incorrect             |
| 500         | Internal server error   | An unexpected error occurred on the server |

## Get User by ID

Retrieve a user by their ID. This endpoint requires authentication.

### Endpoint

```
GET /api/v1/users/{id}
```

### Path Parameters

| Parameter | Type    | Description   |
|-----------|---------|---------------|
| id        | integer | The user's ID |

### Headers

| Header        | Value          | Description                                |
|---------------|----------------|--------------------------------------------|
| Authorization | Bearer {token} | JWT token received from the login endpoint |

### Response

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "employee",
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
}
```

### Example

```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Error Responses

| Status Code | Error Message         | Description                                |
|-------------|-----------------------|--------------------------------------------|
| 400         | Invalid user ID       | The user ID is not a valid number          |
| 401         | Unauthorized          | Missing or invalid JWT token               |
| 404         | User not found        | No user exists with the specified ID       |
| 500         | Internal server error | An unexpected error occurred on the server |

## Get Current User

Retrieve the profile of the currently authenticated user. This endpoint requires authentication.

### Endpoint

```
GET /api/v1/users/me
```

### Headers

| Header        | Value          | Description                                |
|---------------|----------------|--------------------------------------------|
| Authorization | Bearer {token} | JWT token received from the login endpoint |

### Response

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "employee",
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
}
```

### Example

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Error Responses

| Status Code | Error Message         | Description                                |
|-------------|-----------------------|--------------------------------------------|
| 401         | Unauthorized          | Missing or invalid JWT token               |
| 404         | User not found        | The authenticated user no longer exists    |
| 500         | Internal server error | An unexpected error occurred on the server |

## Authentication

Most API endpoints require authentication using a JWT token. To authenticate requests, include the JWT token in the
Authorization header using the Bearer scheme:

```
Authorization: Bearer your_jwt_token
```

You can obtain a JWT token by using the [login endpoint](#user-login).

## JWT Token

The JWT token includes the following claims:

- `user_id`: The ID of the authenticated user
- `email`: The email of the authenticated user
- `role`: The role of the authenticated user
- `exp`: The expiration time (24 hours from token creation)
- `iat`: The token creation time

## Rate Limiting

The API applies rate limiting to prevent abuse. Clients may be restricted to a certain number of requests per minute.

## Error Handling

All errors follow a standard format:

```json
{
  "status": "error",
  "error": {
    "message": "Error message description"
  }
}
```
