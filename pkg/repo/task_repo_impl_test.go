package repo

import (
	"context"
	"testing"
	"time"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/testutil"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupTaskTestRepo sets up a test repository with a test database
func setupTaskTestRepo(t *testing.T) (context.Context, *taskRepoImpl, *userRepoImpl) {
	db := testutil.CreateTestDB(t)
	taskRepo := NewTaskRepoImpl(func(ctx context.Context) *gorm.DB {
		return db.WithContext(ctx)
	})
	userRepo := NewUserRepoImpl(func(ctx context.Context) *gorm.DB {
		return db.WithContext(ctx)
	})
	return t.Context(), taskRepo, userRepo
}

// createTestUser creates a test user with the given details
func createTestUserForTask(t *testing.T, ctx context.Context, r *userRepoImpl, email, name string, role models.UserRole) models.User {
	user, err := r.CreateUser(ctx, models.User{
		Email:        email,
		PasswordHash: "test-password-hash",
		Name:         name,
		Role:         role,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)
	return user
}

// createTestTask creates a test task with the given details
func createTestTask(t *testing.T, ctx context.Context, r *taskRepoImpl, title, description string, employerID models.UserID, assigneeID *models.UserID) models.Task {
	dueDate := time.Now().Add(24 * time.Hour)
	task, err := r.CreateTask(ctx, models.Task{
		Title:       title,
		Description: description,
		Status:      models.TaskStatusPending,
		DueDate:     &dueDate,
		EmployerID:  employerID,
		AssigneeID:  assigneeID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, task.ID)
	return task
}

func Test_taskRepoImpl_CreateTask(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	t.Run("create task without assignee", func(t *testing.T) {
		// Create an employer
		employer := createTestUserForTask(t, ctx, userRepo, "employer1@example.com", "Employer 1", models.UserRoleEmployer)

		// Create a task
		dueDate := time.Now().Add(24 * time.Hour)
		task, err := taskRepo.CreateTask(ctx, models.Task{
			Title:       "Test Task 1",
			Description: "This is a test task",
			Status:      models.TaskStatusPending,
			DueDate:     &dueDate,
			EmployerID:  employer.ID,
		})

		require.NoError(t, err)
		require.NotEmpty(t, task.ID)
		require.Equal(t, "Test Task 1", task.Title)
		require.Equal(t, "This is a test task", task.Description)
		require.Equal(t, models.TaskStatusPending, task.Status)
		require.Equal(t, employer.ID, task.EmployerID)
		require.Nil(t, task.AssigneeID)
	})

	t.Run("create task with assignee", func(t *testing.T) {
		// Create an employer and an employee
		employer := createTestUserForTask(t, ctx, userRepo, "employer2@example.com", "Employer 2", models.UserRoleEmployer)
		employee := createTestUserForTask(t, ctx, userRepo, "employee2@example.com", "Employee 2", models.UserRoleEmployee)

		// Create a task
		dueDate := time.Now().Add(48 * time.Hour)
		task, err := taskRepo.CreateTask(ctx, models.Task{
			Title:       "Test Task 2",
			Description: "This is a test task with assignee",
			Status:      models.TaskStatusPending,
			DueDate:     &dueDate,
			EmployerID:  employer.ID,
			AssigneeID:  &employee.ID,
		})

		require.NoError(t, err)
		require.NotEmpty(t, task.ID)
		require.Equal(t, "Test Task 2", task.Title)
		require.Equal(t, models.TaskStatusPending, task.Status)
		require.Equal(t, employer.ID, task.EmployerID)
		require.NotNil(t, task.AssigneeID)
		require.Equal(t, employee.ID, *task.AssigneeID)
	})
}

func Test_taskRepoImpl_GetTaskByID(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	t.Run("get existing task", func(t *testing.T) {
		// Create users
		employer := createTestUserForTask(t, ctx, userRepo, "employer3@example.com", "Employer 3", models.UserRoleEmployer)
		employee := createTestUserForTask(t, ctx, userRepo, "employee3@example.com", "Employee 3", models.UserRoleEmployee)

		// Create a task
		createdTask := createTestTask(t, ctx, taskRepo, "Task for GetByID", "This is a test task for GetByID", employer.ID, &employee.ID)

		// Get the task by ID
		task, err := taskRepo.GetTaskByID(ctx, createdTask.ID)
		require.NoError(t, err)
		require.NotNil(t, task)
		require.Equal(t, createdTask.ID, task.ID)
		require.Equal(t, "Task for GetByID", task.Title)
		require.Equal(t, employer.ID, task.EmployerID)
		require.NotNil(t, task.AssigneeID)
		require.Equal(t, employee.ID, *task.AssigneeID)
	})

	t.Run("get non-existent task", func(t *testing.T) {
		task, err := taskRepo.GetTaskByID(ctx, models.TaskID(9999))
		require.NoError(t, err)
		require.Nil(t, task)
	})
}

func Test_taskRepoImpl_UpdateTaskStatus(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	t.Run("update status of existing task", func(t *testing.T) {
		// Create users
		employer := createTestUserForTask(t, ctx, userRepo, "employer4@example.com", "Employer 4", models.UserRoleEmployer)

		// Create a task
		task := createTestTask(t, ctx, taskRepo, "Task for UpdateStatus", "This is a test task for UpdateStatus", employer.ID, nil)
		require.Equal(t, models.TaskStatusPending, task.Status)

		// Update the task status to in progress
		updated, err := taskRepo.UpdateTaskStatus(ctx, task.ID, models.TaskStatusInProgress)
		require.NoError(t, err)
		require.True(t, updated)

		// Verify the update
		updatedTask, err := taskRepo.GetTaskByID(ctx, task.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)
		require.Equal(t, models.TaskStatusInProgress, updatedTask.Status)

		// Update the task status to completed
		updated, err = taskRepo.UpdateTaskStatus(ctx, task.ID, models.TaskStatusCompleted)
		require.NoError(t, err)
		require.True(t, updated)

		// Verify the update
		updatedTask, err = taskRepo.GetTaskByID(ctx, task.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)
		require.Equal(t, models.TaskStatusCompleted, updatedTask.Status)
	})

	t.Run("update status of non-existent task", func(t *testing.T) {
		updated, err := taskRepo.UpdateTaskStatus(ctx, models.TaskID(9999), models.TaskStatusInProgress)
		require.NoError(t, err)
		require.False(t, updated)
	})
}

func Test_taskRepoImpl_AssignTask(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	t.Run("assign task to employee", func(t *testing.T) {
		// Create users
		employer := createTestUserForTask(t, ctx, userRepo, "employer5@example.com", "Employer 5", models.UserRoleEmployer)
		employee := createTestUserForTask(t, ctx, userRepo, "employee5@example.com", "Employee 5", models.UserRoleEmployee)

		// Create a task without assignee
		task := createTestTask(t, ctx, taskRepo, "Task for Assign", "This is a test task for Assign", employer.ID, nil)
		require.Nil(t, task.AssigneeID)

		// Assign the task to an employee
		updated, err := taskRepo.AssignTask(ctx, task.ID, employee.ID)
		require.NoError(t, err)
		require.True(t, updated)

		// Verify the assignment
		updatedTask, err := taskRepo.GetTaskByID(ctx, task.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)
		require.NotNil(t, updatedTask.AssigneeID)
		require.Equal(t, employee.ID, *updatedTask.AssigneeID)
	})

	t.Run("reassign task to different employee", func(t *testing.T) {
		// Create users
		employer := createTestUserForTask(t, ctx, userRepo, "employer6@example.com", "Employer 6", models.UserRoleEmployer)
		employee1 := createTestUserForTask(t, ctx, userRepo, "employee6a@example.com", "Employee 6A", models.UserRoleEmployee)
		employee2 := createTestUserForTask(t, ctx, userRepo, "employee6b@example.com", "Employee 6B", models.UserRoleEmployee)

		// Create a task with assignee
		task := createTestTask(t, ctx, taskRepo, "Task for Reassign", "This is a test task for Reassign", employer.ID, &employee1.ID)
		require.NotNil(t, task.AssigneeID)
		require.Equal(t, employee1.ID, *task.AssigneeID)

		// Reassign the task to a different employee
		updated, err := taskRepo.AssignTask(ctx, task.ID, employee2.ID)
		require.NoError(t, err)
		require.True(t, updated)

		// Verify the reassignment
		updatedTask, err := taskRepo.GetTaskByID(ctx, task.ID)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)
		require.NotNil(t, updatedTask.AssigneeID)
		require.Equal(t, employee2.ID, *updatedTask.AssigneeID)
	})

	t.Run("assign non-existent task", func(t *testing.T) {
		employee := createTestUserForTask(t, ctx, userRepo, "employee7@example.com", "Employee 7", models.UserRoleEmployee)

		updated, err := taskRepo.AssignTask(ctx, models.TaskID(9999), employee.ID)
		require.NoError(t, err)
		require.False(t, updated)
	})
}

func Test_taskRepoImpl_GetTasks(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	// Create users
	employer1 := createTestUserForTask(t, ctx, userRepo, "employer8@example.com", "Employer 8", models.UserRoleEmployer)
	employer2 := createTestUserForTask(t, ctx, userRepo, "employer9@example.com", "Employer 9", models.UserRoleEmployer)
	employee1 := createTestUserForTask(t, ctx, userRepo, "employee8@example.com", "Employee 8", models.UserRoleEmployee)
	employee2 := createTestUserForTask(t, ctx, userRepo, "employee9@example.com", "Employee 9", models.UserRoleEmployee)

	// Create tasks
	// Employer 1 tasks
	_ = createTestTask(t, ctx, taskRepo, "Task 1", "Task 1 desc", employer1.ID, &employee1.ID)
	task2 := createTestTask(t, ctx, taskRepo, "Task 2", "Task 2 desc", employer1.ID, &employee1.ID)
	task3 := createTestTask(t, ctx, taskRepo, "Task 3", "Task 3 desc", employer1.ID, &employee2.ID)

	// Employer 2 tasks
	_ = createTestTask(t, ctx, taskRepo, "Task 4", "Task 4 desc", employer2.ID, &employee1.ID)
	task5 := createTestTask(t, ctx, taskRepo, "Task 5", "Task 5 desc", employer2.ID, &employee2.ID)

	// Update some task statuses
	_, err := taskRepo.UpdateTaskStatus(ctx, task2.ID, models.TaskStatusInProgress)
	require.NoError(t, err)

	_, err = taskRepo.UpdateTaskStatus(ctx, task3.ID, models.TaskStatusCompleted)
	require.NoError(t, err)

	_, err = taskRepo.UpdateTaskStatus(ctx, task5.ID, models.TaskStatusInProgress)
	require.NoError(t, err)

	t.Run("get all tasks without filters", func(t *testing.T) {
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{})
		require.NoError(t, err)
		require.Equal(t, 5, count)
		require.Len(t, tasks, 5)
	})

	t.Run("filter by employer", func(t *testing.T) {
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{
			EmployerID: employer1.ID,
		})
		require.NoError(t, err)
		require.Equal(t, 3, count)
		require.Len(t, tasks, 3)

		// Verify all tasks belong to employer1
		for _, task := range tasks {
			require.Equal(t, employer1.ID, task.EmployerID)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{
			Status: models.TaskStatusInProgress,
		})
		require.NoError(t, err)
		require.Equal(t, 2, count)
		require.Len(t, tasks, 2)

		// Verify all tasks are in progress
		for _, task := range tasks {
			require.Equal(t, models.TaskStatusInProgress, task.Status)
		}
	})

	t.Run("filter by employer and status", func(t *testing.T) {
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{
			EmployerID: employer1.ID,
			Status:     models.TaskStatusCompleted,
		})
		require.NoError(t, err)
		require.Equal(t, 1, count)
		require.Len(t, tasks, 1)
		require.Equal(t, task3.ID, tasks[0].ID)
	})

	t.Run("pagination", func(t *testing.T) {
		// Get first page (2 items)
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{
			Limit: 2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, count) // Total count should still be 5
		require.Len(t, tasks, 2)   // But only 2 tasks returned

		// Get second page (2 items)
		count, tasks, err = taskRepo.GetTasks(ctx, GetTasksOptions{
			Offset: 2,
			Limit:  2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, count) // Total count should still be 5
		require.Len(t, tasks, 2)   // But only 2 tasks returned

		// Get third page (1 item)
		count, tasks, err = taskRepo.GetTasks(ctx, GetTasksOptions{
			Offset: 4,
			Limit:  2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, count) // Total count should still be 5
		require.Len(t, tasks, 1)   // But only 1 task returned (the last one)
	})

	t.Run("order by created_at", func(t *testing.T) {
		count, tasks, err := taskRepo.GetTasks(ctx, GetTasksOptions{
			OrderBy: []string{"created_at ASC"},
		})
		require.NoError(t, err)
		require.Equal(t, 5, count)
		require.Len(t, tasks, 5)

		// Verify tasks are ordered by created_at
		for i := 1; i < len(tasks); i++ {
			require.True(t, !tasks[i].CreatedAt.Before(tasks[i-1].CreatedAt))
		}
	})
}

func Test_taskRepoImpl_GetTaskStatistics(t *testing.T) {
	ctx, taskRepo, userRepo := setupTaskTestRepo(t)

	// Create users
	employer := createTestUserForTask(t, ctx, userRepo, "employer10@example.com", "Employer 10", models.UserRoleEmployer)
	employee1 := createTestUserForTask(t, ctx, userRepo, "employee10@example.com", "Employee 10", models.UserRoleEmployee)
	employee2 := createTestUserForTask(t, ctx, userRepo, "employee11@example.com", "Employee 11", models.UserRoleEmployee)

	// Create tasks for employee1
	task1 := createTestTask(t, ctx, taskRepo, "Stat Task 1", "Stat Task 1 desc", employer.ID, &employee1.ID)
	task2 := createTestTask(t, ctx, taskRepo, "Stat Task 2", "Stat Task 2 desc", employer.ID, &employee1.ID)
	task3 := createTestTask(t, ctx, taskRepo, "Stat Task 3", "Stat Task 3 desc", employer.ID, &employee1.ID)

	// Create tasks for employee2
	task4 := createTestTask(t, ctx, taskRepo, "Stat Task 4", "Stat Task 4 desc", employer.ID, &employee2.ID)
	task5 := createTestTask(t, ctx, taskRepo, "Stat Task 5", "Stat Task 5 desc", employer.ID, &employee2.ID)

	// Update task statuses - use the returned tasks for updates
	_, err := taskRepo.UpdateTaskStatus(ctx, task1.ID, models.TaskStatusInProgress)
	require.NoError(t, err)

	_, err = taskRepo.UpdateTaskStatus(ctx, task2.ID, models.TaskStatusCompleted)
	require.NoError(t, err)

	_, err = taskRepo.UpdateTaskStatus(ctx, task4.ID, models.TaskStatusInProgress)
	require.NoError(t, err)

	// Use the tasks in tests to avoid unused variable warnings
	_ = task3
	_ = task5

	t.Run("get statistics for all users", func(t *testing.T) {
		stats, err := taskRepo.GetTaskStatistics(ctx)
		require.NoError(t, err)
		require.Len(t, stats, 2) // Two employees have tasks

		// Check employee1 stats
		emp1Stats, ok := stats[employee1.ID]
		require.True(t, ok)
		require.Equal(t, 3, emp1Stats.TotalTasks)
		require.Equal(t, 1, emp1Stats.Pending)
		require.Equal(t, 1, emp1Stats.InProgress)
		require.Equal(t, 1, emp1Stats.Completed)

		// Check employee2 stats
		emp2Stats, ok := stats[employee2.ID]
		require.True(t, ok)
		require.Equal(t, 2, emp2Stats.TotalTasks)
		require.Equal(t, 1, emp2Stats.Pending)
		require.Equal(t, 1, emp2Stats.InProgress)
		require.Equal(t, 0, emp2Stats.Completed)
	})

	t.Run("get statistics for specific user", func(t *testing.T) {
		stats, err := taskRepo.GetTaskStatistics(ctx, employee1.ID)
		require.NoError(t, err)
		require.Len(t, stats, 1) // Only one employee requested

		// Check employee1 stats
		emp1Stats, ok := stats[employee1.ID]
		require.True(t, ok)
		require.Equal(t, 3, emp1Stats.TotalTasks)
		require.Equal(t, 1, emp1Stats.Pending)
		require.Equal(t, 1, emp1Stats.InProgress)
		require.Equal(t, 1, emp1Stats.Completed)
	})

	t.Run("get statistics for multiple users", func(t *testing.T) {
		stats, err := taskRepo.GetTaskStatistics(ctx, employee1.ID, employee2.ID)
		require.NoError(t, err)
		require.Len(t, stats, 2) // Two employees requested

		// Check stats for both employees
		require.Equal(t, 3, stats[employee1.ID].TotalTasks)
		require.Equal(t, 2, stats[employee2.ID].TotalTasks)
	})

	t.Run("get statistics for user with no tasks", func(t *testing.T) {
		// Create a new employee with no tasks
		employee3 := createTestUserForTask(t, ctx, userRepo, "employee12@example.com", "Employee 12", models.UserRoleEmployee)

		stats, err := taskRepo.GetTaskStatistics(ctx, employee3.ID)
		require.NoError(t, err)
		require.Empty(t, stats) // No stats for this employee
	})
}
