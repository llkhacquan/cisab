package service

import (
	"context"

	"github.com/llkhacquan/cisab/pkg/authctx"
	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/llkhacquan/cisab/pkg/repo"
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

// GetTasks returns all tasks for an employer with filtering, sorting, and pagination
func (s *taskService) GetTasks(ctx context.Context, request GetTasksRequest) (*GetTasksResponse, error) {
	// Check authentication
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}

	// Only employers can access this endpoint
	if authMD.User.Role != models.UserRoleEmployer {
		return nil, NewInvalidInputError("only employers can view all tasks")
	}

	// Build query options
	options := repo.GetTasksOptions{
		Offset: request.Offset,
		Limit:  request.Limit,
	}

	// Add status filter if provided
	if request.Status != "" {
		options.Status = request.Status
	}

	// Add assignee filter if provided
	if request.AssigneeID != nil {
		options.AssigneeID = *request.AssigneeID

		// Verify the assignee exists
		assignee, err := s.userRepo.GetUserByID(ctx, options.AssigneeID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to verify assignee")
		}
		if assignee == nil {
			return nil, NewInvalidInputError("assignee not found")
		}
	}

	// Add sorting options
	if request.SortBy != "" {
		orderBy := request.SortBy

		// Validate sort field
		validSortFields := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"due_date":   true,
			"status":     true,
		}

		if !validSortFields[orderBy] {
			return nil, NewInvalidInputError("sort_by must be one of: created_at, updated_at, due_date, status")
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
		return nil, errors.Wrap(err, "failed to get tasks")
	}

	return &GetTasksResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
	}, nil
}

// AssignTask assigns a task to an employee
func (s *taskService) AssignTask(ctx context.Context, request AssignTaskRequest) (*AssignTaskResponse, error) {
	// Check authentication
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}

	// Only employers can assign tasks
	if authMD.User.Role != models.UserRoleEmployer {
		return nil, NewInvalidInputError("only employers can assign tasks")
	}

	// Get the task
	task, err := s.taskRepo.GetTaskByID(ctx, request.TaskID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task")
	}

	if task == nil {
		return nil, ErrNotFound
	}

	// In this implementation, we allow any employer to assign any task

	// Check if the assignee exists and is an employee
	assignee, err := s.userRepo.GetUserByID(ctx, request.AssigneeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get assignee")
	}
	if assignee == nil {
		return nil, NewInvalidInputError("assignee not found")
	}

	// Verify the assignee is an employee
	if !assignee.IsEmployee() {
		return nil, NewInvalidInputError("tasks can only be assigned to employees")
	}

	// Assign the task
	_, err = s.taskRepo.AssignTask(ctx, request.TaskID, request.AssigneeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to assign task")
	}

	// Retrieve the updated task
	updatedTask, err := s.taskRepo.GetTaskByID(ctx, request.TaskID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get updated task")
	}

	return &AssignTaskResponse{
		Task: *updatedTask,
	}, nil
}

// GetEmployeeTaskSummary returns a summary of task statistics for each employee
func (s *taskService) GetEmployeeTaskSummary(ctx context.Context, request GetEmployeeTaskSummaryRequest) (*GetEmployeeTaskSummaryResponse, error) {
	// Check authentication
	authMD := authctx.Get(ctx)
	if authMD.User.ID == 0 {
		return nil, ErrUnauthorized
	}

	// Only employers can access this endpoint
	if authMD.User.Role != models.UserRoleEmployer {
		return nil, NewInvalidInputError("only employers can view employee task summaries")
	}

	// Get statistics for all users (the repo method returns only employees with tasks)
	statistics, err := s.taskRepo.GetTaskStatistics(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get task statistics")
	}

	// Build the response
	var employeeSummaries []EmployeeSummary
	for userID, stats := range statistics {
		// Get employee details
		employee, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get employee details")
		}

		if employee == nil {
			// This should not happen, but let's handle it anyway
			continue
		}

		// Skip non-employees (just to be safe)
		if employee.Role != models.UserRoleEmployee {
			continue
		}

		employeeSummaries = append(employeeSummaries, EmployeeSummary{
			Employee:   *employee,
			Statistics: stats,
		})
	}

	return &GetEmployeeTaskSummaryResponse{
		Employees: employeeSummaries,
	}, nil
}
