package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/llkhacquan/knovel-assignment/pkg/config"
	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/repo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// userService implements the UserService interface
type userService struct {
	userRepo repo.UserRepo
	jwt      config.JWTConfig
}

// NewUserService creates a new UserService
func NewUserService(userRepo repo.UserRepo, jwt config.JWTConfig) UserService {
	return &userService{
		userRepo: userRepo,
		jwt:      jwt,
	}
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

func (u *userService) GetJWTToken(ctx context.Context, request GetJWTRequest) (*GetJWTResponse, error) {
	// Find user by email
	user, err := u.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find user")
	}
	if user == nil {
		return nil, errors.New("invalid credentials: user not found")
	}

	// Verify password
	if !comparePasswords(user.PasswordHash, request.Password) {
		return nil, errors.New("invalid credentials: password mismatch")
	}

	// Define token expiration time (e.g., 24 hours)
	expirationTime := time.Now().Add(time.Duration(u.jwt.TTLInSecond) * time.Second)

	// Create claims with user information
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(u.jwt.Secret))
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign token")
	}

	// Return response with token and user information
	return &GetJWTResponse{
		Token:       tokenString,
		User:        *user,
		TokenExpiry: expirationTime.Unix(),
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

// comparePasswords compares a hashed password with a plain text password
func comparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
