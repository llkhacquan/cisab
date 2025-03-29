package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/service"
	"github.com/pkg/errors"
)

// GetUserByIDHandler handles GET requests for users
func (s *Server) GetUserByIDHandler(r *http.Request) (interface{}, error) {
	// 1. decode request + basic validation if needed
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.Wrap(err, "invalid user ID")
	}
	var request = service.GetUserByIDRequest{
		ID: models.UserID(id),
	}
	// 2. Call the business logic
	response, err := s.userService.GetUserByID(r.Context(), request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}
	// 3. Return the response
	return response, nil
}

// CreateUserHandler handles POST requests to create a new user.
// curl -X POST http://localhost:8080/users \
// -H "Content-Type: application/json" \
//
//	-d '{
//	 "name": "John Doe",
//	 "email": "john.doe@example.com",
//	 "password": "securepassword",
//	 "role": "employee"
//	}'
func (s *Server) CreateUserHandler(r *http.Request) (interface{}, error) {
	// 1. Decode request
	var createRequest service.CreateUserRequest
	err := ReadJSON(r, &createRequest)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request body")
	}

	// Validate the required fields
	if createRequest.Email == "" || createRequest.Name == "" || createRequest.Password == "" {
		return nil, errors.New("missing required fields")
	}

	// Validate password length
	if len(createRequest.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Validate role
	if createRequest.Role != models.UserRoleEmployee && createRequest.Role != models.UserRoleEmployer {
		return nil, errors.New("invalid role")
	}

	// 2. Call the business logic
	response, err := s.userService.CreateUser(r.Context(), createRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// 3. Return the response
	return response, nil
}

// LoginHandler handles POST requests to authenticate a user and generate a JWT token
// curl -X POST http://localhost:8080/api/v1/login \
// -H "Content-Type: application/json" \
//
//	-d '{
//	  "email": "john.doe@example.com",
//	  "password": "securepassword"
//	}'
func (s *Server) LoginHandler(r *http.Request) (interface{}, error) {
	// 1. Decode request
	var loginRequest service.GetJWTRequest
	err := ReadJSON(r, &loginRequest)
	if err != nil {
		return nil, errors.Wrap(err, "invalid request body")
	}

	// Validate the required fields
	if loginRequest.Email == "" || loginRequest.Password == "" {
		return nil, errors.New("missing required fields")
	}

	// 2. Call the business logic
	response, err := s.userService.GetJWTToken(r.Context(), loginRequest)
	if err != nil {
		return nil, errors.Wrap(err, "authentication failed")
	}

	// 3. Return the response
	return response, nil
}
