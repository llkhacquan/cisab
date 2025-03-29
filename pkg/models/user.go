package models

type UserID int

// User represents a user model
type User struct {
	ID       UserID `json:"id"`
	Username string `json:"username"`
}
