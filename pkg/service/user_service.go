package service

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

// UserService defines the interface for user operations
type UserService interface {
	// GetUsers returns all users
	GetUsers(ctx context.Context, request GetUsersRequest) (GetUsersResponse, error)

	// GetUserByID returns a user by ID
	GetUserByID(ctx context.Context, request GetUserByIDRequest) (*GetUserByIDResponse, error)
}

type GetUsersRequest struct {
	// no fields needed for this request
}
type GetUsersResponse struct {
	Users []models.User `json:"users"`
}

type GetUserByIDRequest struct {
	ID models.UserID
}

type GetUserByIDResponse struct {
	User *models.User `json:"user"`
}
