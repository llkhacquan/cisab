package service

import (
	"context"
	"time"

	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/llkhacquan/cisab/pkg/repo"
)

// TaskService defines the interface for task operations
type TaskService interface {
	// CreateTask creates a new task
	CreateTask(ctx context.Context, request CreateTaskRequest) (*CreateTaskResponse, error)

	// UpdateTaskStatus updates the status of a task
	UpdateTaskStatus(ctx context.Context, request UpdateTaskStatusRequest) (*UpdateTaskStatusResponse, error)

	// AssignTask assigns a task to an employee
	AssignTask(ctx context.Context, request AssignTaskRequest) (*AssignTaskResponse, error)

	// GetAssignedTasks returns tasks assigned to the authenticated employee with filtering and pagination
	GetAssignedTasks(ctx context.Context, request GetAssignedTasksRequest) (*GetAssignedTasksResponse, error)

	// GetTasks returns all tasks for an employer with filtering, sorting, and pagination
	GetTasks(ctx context.Context, request GetTasksRequest) (*GetTasksResponse, error)

	// GetEmployeeTaskSummary returns a summary of task statistics for each employee
	GetEmployeeTaskSummary(ctx context.Context, request GetEmployeeTaskSummaryRequest) (*GetEmployeeTaskSummaryResponse, error)
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

// AssignTaskRequest represents the request to assign a task to an employee
type AssignTaskRequest struct {
	TaskID     models.TaskID `json:"task_id"`
	AssigneeID models.UserID `json:"assignee_id" binding:"required"`
}

// AssignTaskResponse represents the response after assigning a task
type AssignTaskResponse struct {
	Task models.Task `json:"task"`
}

// GetAssignedTasksRequest represents the request for fetching assigned tasks
type GetAssignedTasksRequest struct {
	// Filtering
	Status   models.TaskStatus   `json:"status,omitempty"`   // Filter by specific status (for backward compatibility)
	Statuses []models.TaskStatus `json:"statuses,omitempty"` // Filter by specific statuses

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`    // Field to sort by: "created_at" or "updated_at"
	SortOrder string `json:"sort_order,omitempty"` // Sort order: "asc" or "desc"

	// Pagination
	Limit  int `json:"limit,omitempty"`  // Number of records to return
	Offset int `json:"offset,omitempty"` // Number of records to skip
}

// GetAssignedTasksResponse represents the response for fetching assigned tasks
type GetAssignedTasksResponse struct {
	Tasks      []models.Task `json:"tasks"`
	TotalCount int           `json:"total_count"` // Total number of tasks matching the filters (before pagination)
}

// GetTasksRequest represents the request for fetching tasks (for employers)
type GetTasksRequest struct {
	// Filtering
	Status     models.TaskStatus `json:"status,omitempty"`      // Filter by specific status
	AssigneeID *models.UserID    `json:"assignee_id,omitempty"` // Filter by assignee

	// Sorting
	SortBy    string `json:"sort_by,omitempty"`    // Field to sort by: "created_at", "updated_at", "due_date", or "status"
	SortOrder string `json:"sort_order,omitempty"` // Sort order: "asc" or "desc"

	// Pagination
	Limit  int `json:"limit,omitempty"`  // Number of records to return
	Offset int `json:"offset,omitempty"` // Number of records to skip
}

// GetTasksResponse represents the response for fetching tasks
type GetTasksResponse struct {
	Tasks      []models.Task `json:"tasks"`
	TotalCount int           `json:"total_count"` // Total number of tasks matching the filters (before pagination)
}

// EmployeeSummary represents task statistics for an employee
type EmployeeSummary struct {
	Employee   models.User         `json:"employee"`
	Statistics repo.TaskStatistics `json:"statistics"`
}

// GetEmployeeTaskSummaryRequest represents the request for fetching employee task summaries
type GetEmployeeTaskSummaryRequest struct {
	// Currently empty, but can be extended in the future with filtering and pagination options
}

// GetEmployeeTaskSummaryResponse represents the response for fetching employee task summaries
type GetEmployeeTaskSummaryResponse struct {
	Employees []EmployeeSummary `json:"employees"`
}
