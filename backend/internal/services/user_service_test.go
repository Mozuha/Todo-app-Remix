package services_test

import (
	"context"
	"errors"
	"testing"
	"todo-app/internal/db"
	mock_db "todo-app/internal/db/_mock"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mock_db.NewMockWrappedQuerier(ctrl)
	userService := services.NewUserService(mockQueries)

	uIDStr := "00010203-0405-0607-0809-0a0b0c0d0e0f"
	uIDUuid, _ := utils.StringToUUID(uIDStr)

	t.Run("GetMe", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{UserID: uIDUuid}, nil)

		user, err := userService.GetMe(ctx, uIDUuid)

		require.NoError(t, err)
		assert.Equal(t, uIDUuid, user.UserID)
	})

	t.Run("UpdateUsername", func(t *testing.T) {
		ctx := context.Background()
		req := services.UpdateUsernameRequest{
			Username: "new_username",
		}

		mockQueries.EXPECT().
			UpdateUsername(ctx, db.UpdateUsernameParams{
				Username: req.Username,
				UserID:   uIDUuid,
			}).
			Return(nil)

		err := userService.UpdateUsername(ctx, uIDUuid, req)

		require.NoError(t, err)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			DeleteUser(ctx, uIDUuid).
			Return(db.User{ID: 1, UserID: uIDUuid}, nil)

		err := userService.DeleteUser(ctx, uIDUuid)

		require.NoError(t, err)
	})

	t.Run("GetMe_UserNotFound", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("no rows in result set"))

		user, err := userService.GetMe(ctx, uIDUuid)

		assert.Equal(t, errors.New("no rows in result set"), err)
		assert.Nil(t, user)
	})

	t.Run("GetMe_DBError", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("get failed"))

		user, err := userService.GetMe(ctx, uIDUuid)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("UpdateUsername_DBError", func(t *testing.T) {
		ctx := context.Background()
		req := services.UpdateUsernameRequest{
			Username: "new_username",
		}

		mockQueries.EXPECT().
			UpdateUsername(ctx, db.UpdateUsernameParams{
				Username: req.Username,
				UserID:   uIDUuid,
			}).
			Return(errors.New("update failed"))

		err := userService.UpdateUsername(ctx, uIDUuid, req)

		assert.Error(t, err)
	})

	t.Run("DeleteUser_UserNotFound", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			DeleteUser(ctx, uIDUuid).
			Return(db.User{}, nil)

		err := userService.DeleteUser(ctx, uIDUuid)

		assert.Equal(t, errors.New("no rows in result set"), err)
	})

	t.Run("DeleteUser_DBError", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			DeleteUser(ctx, uIDUuid).
			Return(db.User{}, errors.New("delete failed"))

		err := userService.DeleteUser(ctx, uIDUuid)

		assert.Error(t, err)
	})
}
