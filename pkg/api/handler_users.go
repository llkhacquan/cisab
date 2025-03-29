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
