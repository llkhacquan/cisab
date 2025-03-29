package models

import (
	"time"
)

type UserID int

// UserRole represents the role of a user (employer or employee)
type UserRole string

const (
	// UserRoleEmployer represents an employer user
	UserRoleEmployer UserRole = "employer"
	// UserRoleEmployee represents an employee user
	UserRoleEmployee UserRole = "employee"
)

// User represents a user in the system
type User struct {
	ID           UserID    `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Name         string    `json:"name" gorm:"not null"`
	Role         UserRole  `json:"role" gorm:"type:user_role;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the database table name
func (User) TableName() string {
	return "users"
}

// IsEmployer returns true if the user is an employer
func (u *User) IsEmployer() bool {
	return u.Role == UserRoleEmployer
}

// IsEmployee returns true if the user is an employee
func (u *User) IsEmployee() bool {
	return u.Role == UserRoleEmployee
}
