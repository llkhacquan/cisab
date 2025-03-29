package repo

import (
	"context"

	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var _ TaskRepo = (*taskRepoImpl)(nil)

type taskRepoImpl struct {
	db func(ctx context.Context) *gorm.DB
}

// NewTaskRepoImpl creates a new task repository implementation
func NewTaskRepoImpl(db func(ctx context.Context) *gorm.DB) *taskRepoImpl {
	return &taskRepoImpl{db: db}
}

// GetTaskByID retrieves a task by its ID, return nil if not found
func (r *taskRepoImpl) GetTaskByID(ctx context.Context, id models.TaskID) (*models.Task, error) {
	var task models.Task
	if err := r.db(ctx).Table("tasks").Where("id = ?", id).Limit(1).Scan(&task).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get task by id")
	}
	if task.ID == 0 {
		return nil, nil // Task not found
	}
	return &task, nil
}

// CreateTask creates a new task
func (r *taskRepoImpl) CreateTask(ctx context.Context, task models.Task) (models.Task, error) {
	if err := r.db(ctx).Create(&task).Error; err != nil {
		return models.Task{}, errors.Wrap(err, "failed to create task")
	}
	return task, nil
}

// UpdateTaskStatus updates a task's status
func (r *taskRepoImpl) UpdateTaskStatus(ctx context.Context, id models.TaskID, status models.TaskStatus) (_updated bool, _ error) {
	result := r.db(ctx).Model(&models.Task{}).
		Where("id = ?", id).
		Update("status", status)

	if err := result.Error; err != nil {
		return false, errors.Wrap(err, "failed to update task status")
	}

	return result.RowsAffected > 0, nil
}

// AssignTask assigns a task to an employee (assigneeID)
func (r *taskRepoImpl) AssignTask(ctx context.Context, id models.TaskID, assigneeID models.UserID) (_updated bool, _ error) {
	result := r.db(ctx).Model(&models.Task{}).
		Where("id = ?", id).
		Update("assignee_id", assigneeID)

	if err := result.Error; err != nil {
		return false, errors.Wrap(err, "failed to assign task")
	}

	return result.RowsAffected > 0, nil
}

// GetTasks retrieves all tasks satisfying the given criteria
func (r *taskRepoImpl) GetTasks(ctx context.Context, options GetTasksOptions) (int, []models.Task, error) {
	db := r.db(ctx)

	// Apply filters
	if options.Status != "" {
		db = db.Where("status = ?", options.Status)
	}

	if options.EmployerID != 0 {
		db = db.Where("employer_id = ?", options.EmployerID)
	}

	if options.AssigneeID != 0 {
		db = db.Where("assignee_id = ?", options.AssigneeID)
	}

	// Count total results (before pagination)
	var totalCount int64
	if err := db.Model(&models.Task{}).Count(&totalCount).Error; err != nil {
		return 0, nil, errors.Wrap(err, "failed to count tasks")
	}

	// Apply sorting, pagination, and other options
	for _, orderBy := range options.OrderBy {
		db = db.Order(orderBy)
	}
	if options.Limit > 0 {
		db = db.Limit(options.Limit)
	}
	if options.Offset > 0 {
		db = db.Offset(options.Offset)
	}
	// Fetch results
	var tasks []models.Task
	if err := db.Find(&tasks).Error; err != nil {
		return 0, nil, errors.Wrap(err, "failed to get tasks")
	}

	return int(totalCount), tasks, nil
}

// GetTaskStatistics retrieves the number of tasks for each user
func (r *taskRepoImpl) GetTaskStatistics(ctx context.Context, userIDs ...models.UserID) (map[models.UserID]TaskStatistics, error) {
	result := make(map[models.UserID]TaskStatistics)

	var rows []struct {
		AssigneeID models.UserID
		TotalTasks int
		Pending    int
		InProgress int
		Completed  int
	}
	tx := r.db(ctx).Table("tasks").Select(`
	assignee_id,
	COUNT(*) as total_tasks,
	SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) as pending,
	SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END) as in_progress,
	SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed`)
	if len(userIDs) > 0 {
		tx = tx.Where("assignee_id IN (?)", userIDs)
	}
	if err := tx.Group("assignee_id").Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get task statistics")
	}
	tx = tx.Group("assignee_id")
	if err := tx.Scan(&rows).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get task statistics")
	}
	for _, row := range rows {
		result[row.AssigneeID] = TaskStatistics{
			TotalTasks: row.TotalTasks,
			Pending:    row.Pending,
			InProgress: row.InProgress,
			Completed:  row.Completed,
		}
	}
	return result, nil
}
