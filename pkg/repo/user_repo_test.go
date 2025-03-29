package repo

import (
	"context"
	"testing"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
	"github.com/llkhacquan/knovel-assignment/pkg/testutil"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func Test_userRepoImpl_GetUserByID(t *testing.T) {
	db := testutil.CreateTestDB(t)
	r := NewUserRepoImpl(func(ctx context.Context) *gorm.DB {
		return db.WithContext(ctx)
	})
	ctx := t.Context()
	t.Run("create then get user by id", func(t *testing.T) {
		// Create a test user
		user, err := r.CreateUser(ctx, models.User{
			Username: "test 1",
		})
		require.NoError(t, err)
		require.NotEmpty(t, user.ID)
		require.Equal(t, "test 1", user.Username)
		// Get the user by ID
		user2, err := r.GetUserByID(ctx, user.ID)
		require.NoError(t, err)
		require.NotNil(t, user2)
		require.Equal(t, user.ID, user2.ID)
		require.Equal(t, user.Username, user2.Username)
	})
	t.Run("not found", func(t *testing.T) {
		user, err := r.GetUserByID(ctx, models.UserID(0))
		require.NoError(t, err)
		require.Nil(t, user)
	})
}
