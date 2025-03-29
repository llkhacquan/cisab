package repo

import (
	"context"
	"testing"

	"github.com/llkhacquan/cisab/pkg/models"
	"github.com/llkhacquan/cisab/pkg/testutil"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupTestRepo creates a test database and repository for testing
func setupTestRepo(t *testing.T) (context.Context, UserRepo) {
	db := testutil.CreateTestDB(t)
	r := NewUserRepoImpl(func(ctx context.Context) *gorm.DB {
		return db.WithContext(ctx)
	})
	return t.Context(), r
}

// createTestUser creates a test user with the given details
func createTestUser(t *testing.T, ctx context.Context, r UserRepo, email, name string, role models.UserRole) models.User {
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

func Test_userRepoImpl_CreateUser(t *testing.T) {
	ctx, r := setupTestRepo(t)

	t.Run("create employee user", func(t *testing.T) {
		user, err := r.CreateUser(ctx, models.User{
			Email:        "employee@example.com",
			PasswordHash: "password-hash",
			Name:         "Test Employee",
			Role:         models.UserRoleEmployee,
		})
		require.NoError(t, err)
		require.NotEmpty(t, user.ID)
		require.Equal(t, "employee@example.com", user.Email)
		require.Equal(t, "Test Employee", user.Name)
		require.Equal(t, models.UserRoleEmployee, user.Role)
	})

	t.Run("create employer user", func(t *testing.T) {
		user, err := r.CreateUser(ctx, models.User{
			Email:        "employer@example.com",
			PasswordHash: "password-hash",
			Name:         "Test Employer",
			Role:         models.UserRoleEmployer,
		})
		require.NoError(t, err)
		require.NotEmpty(t, user.ID)
		require.Equal(t, "employer@example.com", user.Email)
		require.Equal(t, "Test Employer", user.Name)
		require.Equal(t, models.UserRoleEmployer, user.Role)
	})
}

func Test_userRepoImpl_GetUserByID(t *testing.T) {
	ctx, r := setupTestRepo(t)

	t.Run("create then get user by id", func(t *testing.T) {
		// Create a test user
		user := createTestUser(t, ctx, r, "user1@example.com", "Test User 1", models.UserRoleEmployee)

		// Get the user by ID
		user2, err := r.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, user2)
		require.Equal(t, user.ID, user2.ID)
		require.Equal(t, user.Email, user2.Email)
		require.Equal(t, user.Name, user2.Name)
		require.Equal(t, user.Role, user2.Role)
	})

	t.Run("not found", func(t *testing.T) {
		user, err := r.GetUserByID(ctx, models.UserID(999))
		require.NoError(t, err)
		require.Nil(t, user)
	})

	t.Run("multiple users", func(t *testing.T) {
		// Create multiple test users
		user1 := createTestUser(t, ctx, r, "multi1@example.com", "Multi User 1", models.UserRoleEmployee)
		user2 := createTestUser(t, ctx, r, "multi2@example.com", "Multi User 2", models.UserRoleEmployer)

		// Get users by ID
		fetchedUser1, err := r.GetUserByID(ctx, user1.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedUser1)
		require.Equal(t, user1.ID, fetchedUser1.ID)
		require.Equal(t, user1.Email, fetchedUser1.Email)

		fetchedUser2, err := r.GetUserByID(ctx, user2.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedUser2)
		require.Equal(t, user2.ID, fetchedUser2.ID)
		require.Equal(t, user2.Email, fetchedUser2.Email)

		// Ensure they're different users
		require.NotEqual(t, fetchedUser1.ID, fetchedUser2.ID)
		require.NotEqual(t, fetchedUser1.Email, fetchedUser2.Email)
	})
}

func Test_userRepoImpl_GetUserByEmail(t *testing.T) {
	ctx, r := setupTestRepo(t)

	t.Run("create then get user by email", func(t *testing.T) {
		// Create a test user
		user := createTestUser(t, ctx, r, "email_test@example.com", "Email Test User", models.UserRoleEmployee)

		// Get the user by email
		user2, err := r.GetUserByEmail(ctx, "email_test@example.com")
		require.NoError(t, err)
		require.NotNil(t, user2)
		require.Equal(t, user.ID, user2.ID)
		require.Equal(t, user.Email, user2.Email)
		require.Equal(t, user.Name, user2.Name)
	})

	t.Run("email not found", func(t *testing.T) {
		user, err := r.GetUserByEmail(ctx, "nonexistent@example.com")
		require.NoError(t, err)
		require.Nil(t, user)
	})

	t.Run("case sensitive email", func(t *testing.T) {
		// Create a test user with specific email
		createTestUser(t, ctx, r, "Case.Sensitive@example.com", "Case Test", models.UserRoleEmployee)

		// Test with different case
		user, err := r.GetUserByEmail(ctx, "case.sensitive@example.com")
		require.NoError(t, err)
		require.Equal(t, "Case.Sensitive@example.com", user.Email)

		// Test with exact case
		user, err = r.GetUserByEmail(ctx, "Case.Sensitive@example.com")
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, "Case.Sensitive@example.com", user.Email)
	})

	t.Run("multiple users with different emails", func(t *testing.T) {
		// Create multiple test users
		user1 := createTestUser(t, ctx, r, "user1@example.com", "User One", models.UserRoleEmployee)
		user2 := createTestUser(t, ctx, r, "user2@example.com", "User Two", models.UserRoleEmployer)

		// Get users by email
		fetchedUser1, err := r.GetUserByEmail(ctx, "user1@example.com")
		require.NoError(t, err)
		require.NotNil(t, fetchedUser1)
		require.Equal(t, user1.ID, fetchedUser1.ID)
		require.Equal(t, "User One", fetchedUser1.Name)

		fetchedUser2, err := r.GetUserByEmail(ctx, "user2@example.com")
		require.NoError(t, err)
		require.NotNil(t, fetchedUser2)
		require.Equal(t, user2.ID, fetchedUser2.ID)
		require.Equal(t, "User Two", fetchedUser2.Name)
	})
}
