package services_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"todo-app/internal/db"
	mock_db "todo-app/internal/db/_mock"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTodoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mock_db.NewMockWrappedQuerier(ctrl)
	todoService := services.NewTodoService(mockQueries)

	uIDStr := "00010203-0405-0607-0809-0a0b0c0d0e0f"
	uIDUuid, _ := utils.StringToUUID(uIDStr)

	t.Run("CreateTodo", func(t *testing.T) {
		ctx := context.Background()
		req := services.CreateTodoRequest{
			Description: "Test todo",
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			CreateTodo(ctx, db.CreateTodoParams{
				UserID:      1,
				Description: req.Description,
			}).
			Return(db.Todo{ID: 1, Description: req.Description}, nil)

		todo, err := todoService.CreateTodo(ctx, uIDUuid, req)

		require.NoError(t, err)
		assert.Equal(t, req.Description, todo.Description)
	})

	t.Run("CreateTodo_UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		req := services.CreateTodoRequest{
			Description: "Test todo",
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		todo, err := todoService.CreateTodo(ctx, uIDUuid, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("CreateTodo_DBError", func(t *testing.T) {
		ctx := context.Background()
		req := services.CreateTodoRequest{
			Description: "Test todo",
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			CreateTodo(ctx, db.CreateTodoParams{
				UserID:      1,
				Description: req.Description,
			}).
			Return(db.Todo{}, errors.New("db error"))

		todo, err := todoService.CreateTodo(ctx, uIDUuid, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("ListTodos", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			ListTodos(ctx, int32(1)).
			Return([]db.Todo{{ID: 1, Description: "Test todo"}}, nil)

		todos, err := todoService.ListTodos(ctx, uIDUuid)

		require.NoError(t, err)
		assert.Len(t, *todos, 1)
	})

	t.Run("ListTodos_UserNotFound", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		todos, err := todoService.ListTodos(ctx, uIDUuid)

		assert.Error(t, err)
		assert.Nil(t, todos)
	})

	t.Run("ListTodos_DBError", func(t *testing.T) {
		ctx := context.Background()

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			ListTodos(ctx, int32(1)).
			Return(nil, errors.New("db error"))

		todos, err := todoService.ListTodos(ctx, uIDUuid)

		assert.Error(t, err)
		assert.Nil(t, todos)
	})

	t.Run("SearchTodos", func(t *testing.T) {
		ctx := context.Background()
		keyword := "Test"

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			SearchTodos(ctx, db.SearchTodosParams{
				UserID:    1,
				ToTsquery: keyword,
			}).
			Return([]db.Todo{{ID: 1, Description: "Test todo"}}, nil)

		todos, err := todoService.SearchTodos(ctx, uIDUuid, keyword)

		require.NoError(t, err)
		assert.Len(t, *todos, 1)
	})

	t.Run("SearchTodos_UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		keyword := "Test"

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		todos, err := todoService.SearchTodos(ctx, uIDUuid, keyword)

		assert.Error(t, err)
		assert.Nil(t, todos)
	})

	t.Run("SearchTodos_DBError", func(t *testing.T) {
		ctx := context.Background()
		keyword := "Test"

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			SearchTodos(ctx, db.SearchTodosParams{
				UserID:    1,
				ToTsquery: keyword,
			}).
			Return(nil, errors.New("db error"))

		todos, err := todoService.SearchTodos(ctx, uIDUuid, keyword)

		assert.Error(t, err)
		assert.Nil(t, todos)
	})

	t.Run("UpdateTodo", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1
		req := services.UpdateTodoRequest{
			Description: "Updated todo",
			Completed:   true,
			Position:    100,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			UpdateTodo(ctx, db.UpdateTodoParams{
				ID:          todoID,
				Description: req.Description,
				Completed:   pgtype.Bool{Bool: req.Completed, Valid: true},
				Position:    pgtype.Numeric{Int: big.NewInt(req.Position), Valid: true},
				UserID:      1,
			}).
			Return(db.Todo{ID: todoID, Description: req.Description, Completed: pgtype.Bool{Bool: req.Completed, Valid: true}}, nil)

		todo, err := todoService.UpdateTodo(ctx, uIDUuid, todoID, req)

		require.NoError(t, err)
		assert.Equal(t, req.Description, todo.Description)
	})

	t.Run("UpdateTodo_UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1
		req := services.UpdateTodoRequest{
			Description: "Updated todo",
			Completed:   true,
			Position:    100,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		todo, err := todoService.UpdateTodo(ctx, uIDUuid, todoID, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("UpdateTodo_DBError", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1
		req := services.UpdateTodoRequest{
			Description: "Updated todo",
			Completed:   true,
			Position:    100,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			UpdateTodo(ctx, db.UpdateTodoParams{
				ID:          todoID,
				Description: req.Description,
				Completed:   pgtype.Bool{Bool: req.Completed, Valid: true},
				Position:    pgtype.Numeric{Int: big.NewInt(req.Position), Valid: true},
				UserID:      1,
			}).
			Return(db.Todo{}, errors.New("db error"))

		todo, err := todoService.UpdateTodo(ctx, uIDUuid, todoID, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	// Let's say a todo C[pos=300] has been moved inbetween A[pos=100] and B[pos=200]
	t.Run("UpdateTodoPosition", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 3
		req := services.UpdateTodoPositionRequest{
			Prevpos: 100,
			Nextpos: 200,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			UpdateTodoPosition(ctx, db.UpdateTodoPositionParams{
				ID:      todoID,
				UserID:  1,
				Prevpos: pgtype.Numeric{Int: big.NewInt(req.Prevpos), Valid: true},
				Nextpos: pgtype.Numeric{Int: big.NewInt(req.Nextpos), Valid: true},
			}).
			Return(db.Todo{ID: todoID, Position: pgtype.Numeric{Int: big.NewInt(150), Valid: true}}, nil)

		todo, err := todoService.UpdateTodoPosition(ctx, uIDUuid, todoID, req)

		require.NoError(t, err)
		assert.Equal(t, todoID, todo.ID)
	})

	t.Run("UpdateTodoPosition_UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 3
		req := services.UpdateTodoPositionRequest{
			Prevpos: 100,
			Nextpos: 200,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		todo, err := todoService.UpdateTodoPosition(ctx, uIDUuid, todoID, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("UpdateTodoPosition_DBError", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 3
		req := services.UpdateTodoPositionRequest{
			Prevpos: 100,
			Nextpos: 200,
		}

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			UpdateTodoPosition(ctx, db.UpdateTodoPositionParams{
				ID:      todoID,
				UserID:  1,
				Prevpos: pgtype.Numeric{Int: big.NewInt(req.Prevpos), Valid: true},
				Nextpos: pgtype.Numeric{Int: big.NewInt(req.Nextpos), Valid: true},
			}).
			Return(db.Todo{}, errors.New("db error"))

		todo, err := todoService.UpdateTodoPosition(ctx, uIDUuid, todoID, req)

		assert.Error(t, err)
		assert.Nil(t, todo)
	})

	t.Run("DeleteTodo", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			DeleteTodo(ctx, db.DeleteTodoParams{
				ID:     todoID,
				UserID: 1,
			}).
			Return(db.Todo{ID: todoID}, nil)

		err := todoService.DeleteTodo(ctx, uIDUuid, todoID)

		require.NoError(t, err)
	})

	t.Run("DeleteTodo_UserNotFound", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{}, errors.New("user not found"))

		err := todoService.DeleteTodo(ctx, uIDUuid, todoID)

		assert.Error(t, err)
	})

	t.Run("DeleteTodo_TodoNotFound", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1000

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			DeleteTodo(ctx, db.DeleteTodoParams{
				ID:     todoID,
				UserID: 1,
			}).
			Return(db.Todo{}, nil)

		err := todoService.DeleteTodo(ctx, uIDUuid, todoID)

		assert.Equal(t, utils.ErrNoRowsMatchedSQLC, err)
	})

	t.Run("DeleteTodo_DBError", func(t *testing.T) {
		ctx := context.Background()
		var todoID int32 = 1

		mockQueries.EXPECT().
			GetUserByUserID(ctx, uIDUuid).
			Return(db.User{ID: 1}, nil)

		mockQueries.EXPECT().
			DeleteTodo(ctx, db.DeleteTodoParams{
				ID:     todoID,
				UserID: 1,
			}).
			Return(db.Todo{}, errors.New("db error"))

		err := todoService.DeleteTodo(ctx, uIDUuid, todoID)

		assert.Error(t, err)
	})
}
