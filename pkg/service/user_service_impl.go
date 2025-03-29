package service

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/repo"
	"github.com/pkg/errors"
)

// userService implements the UserService interface
type userService struct {
	userRepo repo.UserRepo
}

func (u *userService) GetUsers(ctx context.Context, request GetUsersRequest) (GetUsersResponse, error) {
	// This is a placeholder implementation
	// In a real application, you would fetch users from a database with pagination
	// For now, we'll return hardcoded data
	return GetUsersResponse{
		Users: []models.User{
			{ID: 1, Name: "johndoe", Email: "john@example.com", Role: models.UserRoleEmployee},
			{ID: 2, Name: "janedoe", Email: "jane@example.com", Role: models.UserRoleEmployer},
			{ID: 3, Name: "bobsmith", Email: "bob@example.com", Role: models.UserRoleEmployee},
		},
	}, nil
}

func (u *userService) GetUserByID(ctx context.Context, request GetUserByIDRequest) (*GetUserByIDResponse, error) {
	user, err := u.userRepo.GetUserByID(ctx, request.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &GetUserByIDResponse{
		User: user,
	}, nil
}

// NewUserService creates a new UserService
func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
