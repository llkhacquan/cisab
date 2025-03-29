package service

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

// UserService defines the interface for user operations
type UserService interface {
	// GetUsers returns all users
	GetUsers(ctx context.Context) ([]models.User, error)

	// GetUserByID returns a user by ID
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

// userService implements the UserService interface
type userService struct {
	// In a real application, this would have dependencies like a database repository
}

// NewUserService creates a new UserService
func NewUserService() UserService {
	return &userService{}
}

// GetUsers returns all users
func (s *userService) GetUsers(ctx context.Context) ([]models.User, error) {
	// In a real application, this would fetch users from a database
	// For now, we'll return mock data
	return []models.User{
		{ID: 1, Username: "johndoe", Email: "john@example.com", Role: "admin"},
		{ID: 2, Username: "janedoe", Email: "jane@example.com", Role: "user"},
		{ID: 3, Username: "bobsmith", Email: "bob@example.com", Role: "user"},
	}, nil
}

// GetUserByID returns a user by ID
func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	users, err := s.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, nil // User not found
}
