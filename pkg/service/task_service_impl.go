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
