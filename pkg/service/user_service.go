package service

import (
	"context"

	"github.com/llkhacquan/cisab/pkg/models"
)

// UserService defines the interface for user operations
type UserService interface {
	// GetUserByID returns a user by ID
	GetUserByID(ctx context.Context, request GetUserByIDRequest) (*GetUserByIDResponse, error)

	// GetMe returns the current authenticated user
	GetMe(ctx context.Context) (*GetMeResponse, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, error)

	// GetJWTToken generates and returns a JWT token for authentication
	GetJWTToken(ctx context.Context, request GetJWTRequest) (*GetJWTResponse, error)

	// GetUsers returns all users (only accessible by employers)
	GetUsers(ctx context.Context, request GetUsersRequest) (*GetUsersResponse, error)
}

type GetUserByIDRequest struct {
	ID models.UserID
}

type GetUserByIDResponse struct {
	User *models.User `json:"user"`
}

type GetMeResponse struct {
	User *models.User `json:"user"`
}

type CreateUserRequest struct {
	Name     string          `json:"name" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     models.UserRole `json:"role" binding:"required"`
}

type CreateUserResponse struct {
	User models.User `json:"user"`
}

type GetJWTRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GetJWTResponse struct {
	Token       string      `json:"token"`
	User        models.User `json:"user"`
	TokenExpiry int64       `json:"token_expiry"`
}

type GetUsersRequest struct {
	// No filters needed for this simple implementation
}

type GetUsersResponse struct {
	Users []models.User `json:"users"`
}
