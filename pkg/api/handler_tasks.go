package api

import (
	"net/http"

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
