package service

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/repo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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

func (u *userService) CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, error) {
	// we should check if the user already exists (by email). But let it be for now, the DB will handle it
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}
	// Create the user
	newUser := models.User{
		Email:        request.Email,
		PasswordHash: hashedPassword,
		Name:         request.Name,
		Role:         request.Role,
	}
	createdUser, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}
	return &CreateUserResponse{
		User: createdUser,
	}, nil
}

// hashPassword is a simple password hashing function using bcrypt with default cost
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// NewUserService creates a new UserService
func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
