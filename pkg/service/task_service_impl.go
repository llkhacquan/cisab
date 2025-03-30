package service

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/authctx"
	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/repo"
	"github.com/pkg/errors"
)

// taskService implements the TaskService interface
type taskService struct {
	taskRepo repo.TaskRepo
	userRepo repo.UserRepo
}

// NewTaskService creates a new TaskService
func NewTaskService(taskRepo repo.TaskRepo, userRepo repo.UserRepo) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// CreateTask creates a new task
func (s *taskService) CreateTask(ctx context.Context, request CreateTaskRequest) (*CreateTaskResponse, error) {
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}
	if authMD.User.Role != models.UserRoleEmployer {
		return nil, NewInvalidInputError("only employers can create tasks")
	}

	status := models.TaskStatusPending
	if request.Status != "" {
		switch models.TaskStatus(request.Status) {
		case models.TaskStatusPending, models.TaskStatusInProgress, models.TaskStatusCompleted:
			status = models.TaskStatus(request.Status)
		default:
			return nil, NewInvalidInputError("invalid task status")
		}
	}

	// Check if the assignee exists if provided
	var assigneeID *models.UserID
	if request.AssigneeID != nil {
		id := models.UserID(*request.AssigneeID)
		assignee, err := s.userRepo.GetUserByID(ctx, id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get assignee")
		}
		if assignee == nil {
			return nil, NewInvalidInputError("assignee not found")
		}
		// Verify the assignee is an employee
		if !assignee.IsEmployee() {
			return nil, errors.New("task can only be assigned to employees")
		}
		assigneeID = &id
	}

	// Create the task
	task := models.Task{
		Title:       request.Title,
		Description: request.Description,
		Status:      status,
		DueDate:     request.DueDate,
		EmployerID:  authMD.User.ID,
		AssigneeID:  assigneeID,
	}

	createdTask, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create task")
	}

	return &CreateTaskResponse{
		Task: createdTask,
	}, nil
}

// UpdateTaskStatus updates the status of a task
func (s *taskService) UpdateTaskStatus(ctx context.Context, request UpdateTaskStatusRequest) (*UpdateTaskStatusResponse, error) {
	// Check authentication
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}

	// Get the task
	task, err := s.taskRepo.GetTaskByID(ctx, request.TaskID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task")
	}

	if task == nil {
		return nil, ErrNotFound
	}

	// Verify the user is allowed to update the task status
	// Only employees assigned to the task can update its status
	if authMD.User.Role == models.UserRoleEmployee {
		// Make sure task has assignee and it's the current user
		if task.AssigneeID == nil || *task.AssigneeID != authMD.User.ID {
			return nil, NewInvalidInputError("you can only update tasks assigned to you")
		}
	} else if authMD.User.Role == models.UserRoleEmployer {
		// Employers can only update tasks they created
		if task.EmployerID != authMD.User.ID {
			return nil, NewInvalidInputError("you can only update tasks you created")
		}
	}

	// Validate the new status
	switch request.Status {
	case models.TaskStatusPending, models.TaskStatusInProgress, models.TaskStatusCompleted:
		// Valid status
	default:
		return nil, NewInvalidInputError("invalid task status")
	}

	// Update the task status
	_, err = s.taskRepo.UpdateTaskStatus(ctx, task.ID, request.Status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update task status")
	}

	// Note: We're ignoring the 'updated' boolean return value since we always fetch the task afterward

	// Retrieve the updated task
	updatedTask, err := s.taskRepo.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated task")
	}

	return &UpdateTaskStatusResponse{
		Task: *updatedTask,
	}, nil
}

// GetAssignedTasks returns tasks assigned to the authenticated employee with filtering and pagination
func (s *taskService) GetAssignedTasks(ctx context.Context, request GetAssignedTasksRequest) (*GetAssignedTasksResponse, error) {
	// Check authentication
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}

	// Employee can only view their assigned tasks
	if authMD.User.Role != models.UserRoleEmployee {
		return nil, NewInvalidInputError("only employees can view their assigned tasks")
	}

	// Build query options
	options := repo.GetTasksOptions{
		AssigneeID: authMD.User.ID,
		Offset:     request.Offset,
		Limit:      request.Limit,
	}

	// Add status filter if provided
	if len(request.Statuses) > 0 {
		// Note: The current repo implementation only supports filtering by a single status,
		// so we're using the first status in the list
		options.Status = request.Statuses[0]
	} else if request.Status != "" {
		// For backward compatibility
		options.Status = request.Status
	}

	// Add sorting options
	if request.SortBy != "" {
		orderBy := request.SortBy

		// Validate sort field
		if orderBy != "created_at" && orderBy != "updated_at" {
			return nil, NewInvalidInputError("sort_by must be 'created_at' or 'updated_at'")
		}

		// Add sort order
		if request.SortOrder != "" {
			if request.SortOrder != "asc" && request.SortOrder != "desc" {
				return nil, NewInvalidInputError("sort_order must be 'asc' or 'desc'")
			}
			orderBy += " " + request.SortOrder
		} else {
			// Default to descending order
			orderBy += " DESC"
		}

		options.OrderBy = []string{orderBy}
	} else {
		// Default sort order: created_at descending (newest first)
		options.OrderBy = []string{"created_at DESC"}
	}

	// Get tasks from repository
	totalCount, tasks, err := s.taskRepo.GetTasks(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get assigned tasks")
	}

	return &GetAssignedTasksResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
	}, nil
}
