package repo

import (
	"context"

	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ UserRepo = (*userRepoImpl)(nil)

type userRepoImpl struct {
	db func(ctx context.Context) *gorm.DB
}

func NewUserRepoImpl(db func(ctx context.Context) *gorm.DB) *userRepoImpl {
	return &userRepoImpl{db: db}
}

func (u userRepoImpl) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	if err := u.db(ctx).Create(&user).Error; err != nil {
		return models.User{}, errors.Wrap(err, "failed to create user")
	}
	return user, nil
}

func (u userRepoImpl) GetUserByID(ctx context.Context, id models.UserID) (*models.User, error) {
	var user models.User
	if err := u.db(ctx).Table("users").Where("id = ?", id).Limit(1).Scan(&user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user by id")
	}
	if user.ID == 0 {
		return nil, nil // user not found
	}
	return &user, nil
}

func (u userRepoImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := u.db(ctx).Table("users").Where("email = ?", email).Limit(1).Scan(&user).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	if user.ID == 0 {
		return nil, nil // user not found
	}
	return &user, nil
}

func (u userRepoImpl) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := u.db(ctx).Table("users").Order("id asc").Scan(&users).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}
	return users, nil
}
