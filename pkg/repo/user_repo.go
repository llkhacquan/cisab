package repo

import (
	"context"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

type UserRepo interface {
	// GetUserByID retrieves a user by their ID, return nil if not found
	GetUserByID(ctx context.Context, id models.UserID) (*models.User, error)
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	// GetUserByEmail retrieves a user by their email, return nil if not found.
	// This is useful for login or registration processes.
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
