package service

import (
	"context"
	"time"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

// TaskService defines the interface for task operations
type TaskService interface {
	// CreateTask creates a new task
	CreateTask(ctx context.Context, request CreateTaskRequest) (*CreateTaskResponse, error)

	// UpdateTaskStatus updates the status of a task
	UpdateTaskStatus(ctx context.Context, request UpdateTaskStatusRequest) (*UpdateTaskStatusResponse, error)
}

// CreateTaskRequest represents the request to create a new task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssigneeID  *int       `json:"assignee_id,omitempty"`
}

// CreateTaskResponse represents the response after creating a task
type CreateTaskResponse struct {
	Task models.Task `json:"task"`
}

// UpdateTaskStatusRequest represents the request to update a task's status
type UpdateTaskStatusRequest struct {
	TaskID models.TaskID     `json:"task_id"`
	Status models.TaskStatus `json:"status" binding:"required"`
}

// UpdateTaskStatusResponse represents the response after updating a task's status
type UpdateTaskStatusResponse struct {
	Task models.Task `json:"task"`
}
