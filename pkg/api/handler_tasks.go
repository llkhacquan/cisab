package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/service"
	"github.com/pkg/errors"
)

// CreateTaskHandler handles POST requests to create a new task
// curl -X POST http://localhost:8080/api/v1/tasks \
// -H "Content-Type: application/json" \
// -H "Authorization: Bearer your_jwt_token" \
//
//	-d '{
//	  "title": "Task Title",
//	  "description": "Task Description",
//	  "due_date": "2023-04-15T00:00:00Z"
//	}'
func (s *Server) CreateTaskHandler(r *http.Request) (interface{}, error) {
	// 1. Decode request
	var createRequest service.CreateTaskRequest
	err := ReadJSON(r, &createRequest)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request body")
	}

	// Validate the required fields
	if createRequest.Title == "" {
		return nil, errors.New("title is required")
	}

	// 2. Call the business logic
	response, err := s.taskService.CreateTask(r.Context(), createRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create task")
	}

	// 3. Return the response
	return response, nil
}

// UpdateTaskStatusHandler handles PATCH requests to update a task's status
// curl -X PATCH http://localhost:8080/api/v1/tasks/{id}/status \
// -H "Content-Type: application/json" \
// -H "Authorization: Bearer your_jwt_token" \
//
//	-d '{
//	  "status": "completed"
//	}'
func (s *Server) UpdateTaskStatusHandler(r *http.Request) (interface{}, error) {
	// 1. Extract task ID from URL
	vars := mux.Vars(r)
	taskIDStr := vars["id"]
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid task ID")
	}

	// 2. Decode request body
	var updateRequest struct {
		Status string `json:"status"`
	}
	if err := ReadJSON(r, &updateRequest); err != nil {
		return nil, errors.Wrap(err, "invalid request body")
	}

	// Validate status
	if updateRequest.Status == "" {
		return nil, errors.New("status is required")
	}

	// 3. Create service request
	serviceRequest := service.UpdateTaskStatusRequest{
		TaskID: models.TaskID(taskID),
		Status: models.TaskStatus(updateRequest.Status),
	}

	// 4. Call the business logic
	response, err := s.taskService.UpdateTaskStatus(r.Context(), serviceRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update task status")
	}

	// 5. Return the response
	return response, nil
}

// GetAssignedTasksHandler handles GET requests to retrieve tasks assigned to the authenticated employee
// curl -X GET http://localhost:8080/api/v1/tasks/assigned \
// -H "Authorization: Bearer your_jwt_token" \
// "?status=pending&sort_by=created_at&sort_order=desc&limit=10&offset=0"
func (s *Server) GetAssignedTasksHandler(r *http.Request) (interface{}, error) {
	// Parse query parameters
	query := r.URL.Query()

	// Build the request
	request := service.GetAssignedTasksRequest{}

	// Parse statuses filter
	if statusStr := query.Get("status"); statusStr != "" {
		request.Statuses = []models.TaskStatus{models.TaskStatus(statusStr)}
	}

	// Parse sorting options
	request.SortBy = query.Get("sort_by")
	request.SortOrder = query.Get("sort_order")

	// Parse pagination options
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid limit parameter")
		}
		request.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid offset parameter")
		}
		request.Offset = offset
	}

	// Call the business logic
	response, err := s.taskService.GetAssignedTasks(r.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get assigned tasks")
	}

	// Return the response
	return response, nil
}

// GetEmployeeTaskSummaryHandler handles GET requests to retrieve task statistics for all employees
// curl -X GET http://localhost:8080/api/v1/employee-summary \
// -H "Authorization: Bearer your_jwt_token"
func (s *Server) GetEmployeeTaskSummaryHandler(r *http.Request) (interface{}, error) {
	// Create an empty request
	request := service.GetEmployeeTaskSummaryRequest{}

	// Call the business logic
	response, err := s.taskService.GetEmployeeTaskSummary(r.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get employee task summary")
	}

	// Return the response
	return response, nil
}

// GetTasksHandler handles GET requests to retrieve tasks with filtering and sorting
// curl -X GET http://localhost:8080/api/v1/tasks \
// -H "Authorization: Bearer your_jwt_token" \
// "?status=pending&assignee_id=1&sort_by=due_date&sort_order=asc&limit=10&offset=0"
func (s *Server) GetTasksHandler(r *http.Request) (interface{}, error) {
	// Parse query parameters
	query := r.URL.Query()

	// Build the request
	request := service.GetTasksRequest{}

	// Parse status filter
	if status := query.Get("status"); status != "" {
		request.Status = models.TaskStatus(status)
	}

	// Parse assignee filter
	if assigneeIDStr := query.Get("assignee_id"); assigneeIDStr != "" {
		assigneeID, err := strconv.Atoi(assigneeIDStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid assignee_id parameter")
		}
		id := models.UserID(assigneeID)
		request.AssigneeID = &id
	}

	// Parse sorting options
	request.SortBy = query.Get("sort_by")
	request.SortOrder = query.Get("sort_order")

	// Parse pagination options
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid limit parameter")
		}
		request.Limit = limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid offset parameter")
		}
		request.Offset = offset
	}

	// Call the business logic
	response, err := s.taskService.GetTasks(r.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tasks")
	}

	// Return the response
	return response, nil
}
