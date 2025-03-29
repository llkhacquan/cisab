package models

import (
	"time"
)

type TaskID int

// TaskStatus represents the status of a task
type TaskStatus string

const (
	// TaskStatusPending represents a task that has not been started
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusInProgress represents a task that is in progress
	TaskStatusInProgress TaskStatus = "in_progress"
	// TaskStatusCompleted represents a task that has been completed
	TaskStatusCompleted TaskStatus = "completed"
)

// Task represents a task in the system
type Task struct {
	ID          TaskID     `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description" gorm:"type:text"`
	Status      TaskStatus `json:"status" gorm:"not null;default:'pending'"`
	DueDate     *time.Time `json:"due_date,omitempty" gorm:"index"`
	EmployerID  UserID     `json:"employer_id" gorm:"not null;index"`
	AssigneeID  *UserID    `json:"assignee_id,omitempty" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Define relationships (not stored in database)
	Employer *User `json:"employer,omitempty" gorm:"foreignKey:EmployerID"`
	Assignee *User `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
}

// TableName specifies the database table name
func (Task) TableName() string {
	return "tasks"
}

// IsAssigned returns true if the task has been assigned to an employee
func (t *Task) IsAssigned() bool {
	return t.AssigneeID != nil
}

// IsPending returns true if the task is pending
func (t *Task) IsPending() bool {
	return t.Status == TaskStatusPending
}

// IsInProgress returns true if the task is in progress
func (t *Task) IsInProgress() bool {
	return t.Status == TaskStatusInProgress
}

// IsCompleted returns true if the task is completed
func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// IsOverdue returns true if the task has a due date and it has passed
func (t *Task) IsOverdue() bool {
	return t.DueDate != nil && time.Now().After(*t.DueDate) && !t.IsCompleted()
}
