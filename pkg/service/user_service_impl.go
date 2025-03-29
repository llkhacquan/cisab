package service

import (
	"context"
)

// userService implements the UserService interface
type userService struct {
	// In a real application, this would have dependencies like a database repository
}

func (u userService) GetUsers(ctx context.Context, request GetUsersRequest) (GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u userService) GetUserByID(ctx context.Context, request GetUserByIDRequest) (*GetUserByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewUserService creates a new UserService
func NewUserService() UserService {
	return &userService{}
}
