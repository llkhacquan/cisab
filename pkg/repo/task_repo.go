package repo

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

type TaskRepo interface {
	// GetTaskByID retrieves a task by its ID, return nil if not found
	GetTaskByID(ctx context.Context, id models.TaskID) (*models.Task, error)
	// CreateTask creates a new task.
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
	// UpdateTaskStatus updates a task's status.
	// In real-world applications, this would likely involve more complex logic, but for this assignment,
	// we'll keep it simple by just updating the status field only
	UpdateTaskStatus(ctx context.Context, id models.TaskID, status models.TaskStatus) (_updated bool, _ error)
	// AssignTask assigns a task to an employee (assigneeID)
	AssignTask(ctx context.Context, id models.TaskID, assigneeID models.UserID) (_updated bool, _ error)

	// GetTasks retrieves all tasks satisfying the given criteria.
	GetTasks(ctx context.Context, options GetTasksOptions) (_total int, _ []models.Task, _ error)
	// GetTaskStatistics retrieves the number of tasks for each user.
	// If no userIDs are provided, it retrieves statistics for all users.
	GetTaskStatistics(ctx context.Context, userIDs ...models.UserID) (map[models.UserID]TaskStatistics, error)
}

type GetTasksOptions struct {
	// Filter by status
	Status     models.TaskStatus // if not set, all statuses are included
	EmployerID models.UserID     // if not set, all employers are included

	// this is sql-like syntax, it might not be supported by all databases,
	// but for the sake of this assignment, we'll assume it's supported.
	OrderBy []string // the input should be in the format of "column_name ASC|DESC"
	Offset  int      // offset for pagination
	Limit   int      // limit for pagination
}

type TaskStatistics struct {
	TotalTasks int `json:"total_tasks"`
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
}
