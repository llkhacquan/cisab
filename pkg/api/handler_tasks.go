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
