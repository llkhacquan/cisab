# Task Service API Documentation

This document provides detailed information about the Task Service API endpoints, including request/response formats and
curl examples.

## Table of Contents

1. [Create Task](#create-task)
2. [Update Task Status](#update-task-status)

## Create Task

Create a new task. This endpoint requires authentication and can only be used by employers.

### Endpoint

```
POST /api/v1/tasks
```

### Request Body

```json
{
  "title": "string",
  "description": "string",
  "status": "pending|in_progress|completed",
  "due_date": "ISO 8601 datetime string",
  "assignee_id": "integer (optional)"
}
```

#### Fields

| Field       | Type    | Required | Description                                                |
|-------------|---------|----------|------------------------------------------------------------|
| title       | string  | Yes      | Task title                                                 |
| description | string  | No       | Detailed description of the task                           |
| status      | string  | No       | Task status (defaults to "pending" if not provided)        |
| due_date    | string  | No       | Due date in ISO 8601 format (e.g., "2023-04-15T00:00:00Z") |
| assignee_id | integer | No       | ID of the employee to assign the task to                   |

### Response

```json
{
  "status": "success",
  "data": {
    "task": {
      "id": 1,
      "title": "Task Title",
      "description": "Task Description",
      "status": "pending",
      "due_date": "2023-04-15T00:00:00Z",
      "employer_id": 2,
      "assignee_id": 3,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
}
```

### Example

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "title": "Implement user authentication",
    "description": "Add JWT-based authentication to the API",
    "due_date": "2023-04-15T00:00:00Z",
    "assignee_id": 3
  }'
```

### Error Responses

| Status Code | Error Message                   | Description                                             |
|-------------|---------------------------------|---------------------------------------------------------|
| 400         | Missing required fields         | One or more required fields are missing                 |
| 400         | Invalid task status             | Status must be "pending", "in_progress", or "completed" |
| 400         | Invalid assignee                | The assignee specified does not exist                   |
| 400         | Assignee must be an employee    | Only employees can be assigned to tasks                 |
| 401         | Unauthorized                    | Missing or invalid JWT token                            |
| 403         | Only employers can create tasks | The authenticated user is not an employer               |
| 500         | Internal server error           | An unexpected error occurred on the server              |

## Update Task Status

Update the status of a task. This endpoint requires authentication and can be used by:

- Employees: Can only update the status of tasks assigned to them
- Employers: Can only update the status of tasks they created

### Endpoint

```
PATCH /api/v1/tasks/{id}/status
```

### URL Parameters

| Parameter | Type    | Description                  |
|-----------|---------|------------------------------|
| id        | integer | The ID of the task to update |

### Request Body

```json
{
  "status": "pending|in_progress|completed"
}
```

#### Fields

| Field  | Type   | Required | Description                                                      |
|--------|--------|----------|------------------------------------------------------------------|
| status | string | Yes      | The new status for the task (pending, in_progress, or completed) |

### Response

```json
{
  "status": "success",
  "data": {
    "task": {
      "id": 1,
      "title": "Task Title",
      "description": "Task Description",
      "status": "completed",
      "due_date": "2023-04-15T00:00:00Z",
      "employer_id": 2,
      "assignee_id": 3,
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T14:30:00Z"
    }
  }
}
```

### Example

```bash
curl -X PATCH http://localhost:8080/api/v1/tasks/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "status": "completed"
  }'
```

### Error Responses

| Status Code | Error Message                             | Description                                                  |
|-------------|-------------------------------------------|--------------------------------------------------------------|
| 400         | Missing required fields                   | Status field is missing                                      |
| 400         | Invalid task status                       | Status must be "pending", "in_progress", or "completed"      |
| 400         | You can only update tasks assigned to you | The employee is trying to update a task not assigned to them |
| 400         | You can only update tasks you created     | The employer is trying to update a task they didn't create   |
| 401         | Unauthorized                              | Missing or invalid JWT token                                 |
| 404         | Not found                                 | The specified task does not exist                            |
| 500         | Internal server error                     | An unexpected error occurred on the server                   |
