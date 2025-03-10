package services

import (
	"context"
	"math/big"
	"todo-app/internal/db"
	"todo-app/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type TodoService struct {
	SqlClient db.WrappedQuerier
}

type CreateTodoRequest struct {
	Description string `json:"description" binding:"required"`
}

type UpdateTodoRequest struct {
	Description string `json:"description" binding:"required"`
	Completed   bool   `json:"completed"`
	Position    int64  `json:"position" binding:"required"`
}

type UpdateTodoPositionRequest struct {
	Prevpos int64 `json:"prev_pos" binding:"required"`
	Nextpos int64 `json:"next_pos" binding:"required"`
}

func NewTodoService(sqlClient db.WrappedQuerier) *TodoService {
	return &TodoService{SqlClient: sqlClient}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID pgtype.UUID, req CreateTodoRequest) (*db.Todo, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInvalidUID
	}

	todo, err := s.SqlClient.CreateTodo(ctx, db.CreateTodoParams{
		UserID:      user.ID,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoService) ListTodos(ctx context.Context, userID pgtype.UUID) (*[]db.Todo, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInvalidUID
	}

	todos, err := s.SqlClient.ListTodos(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &todos, nil
}

func (s *TodoService) SearchTodos(ctx context.Context, userID pgtype.UUID, keyword string) (*[]db.Todo, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInvalidUID
	}

	todos, err := s.SqlClient.SearchTodos(ctx, db.SearchTodosParams{
		UserID:    user.ID,
		ToTsquery: keyword,
	})
	if err != nil {
		return nil, err
	}

	return &todos, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, userID pgtype.UUID, todoID int32, req UpdateTodoRequest) (*db.Todo, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInvalidUID
	}

	todo, err := s.SqlClient.UpdateTodo(ctx, db.UpdateTodoParams{
		ID:          todoID,
		Description: req.Description,
		Completed:   pgtype.Bool{Bool: req.Completed, Valid: true},
		Position:    pgtype.Numeric{Int: big.NewInt(req.Position), Valid: true},
		UserID:      user.ID,
	})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoService) UpdateTodoPosition(ctx context.Context, userID pgtype.UUID, todoID int32, req UpdateTodoPositionRequest) (*db.Todo, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, utils.ErrInvalidUID
	}

	todo, err := s.SqlClient.UpdateTodoPosition(ctx, db.UpdateTodoPositionParams{
		ID:      todoID,
		UserID:  user.ID,
		Prevpos: pgtype.Numeric{Int: big.NewInt(req.Prevpos), Valid: true},
		Nextpos: pgtype.Numeric{Int: big.NewInt(req.Nextpos), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s *TodoService) DeleteTodo(ctx context.Context, userID pgtype.UUID, todoID int32) error {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return utils.ErrInvalidUID
	}

	todo, err := s.SqlClient.DeleteTodo(ctx, db.DeleteTodoParams{
		ID:     todoID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	} else if todo.ID == 0 {
		return utils.ErrNoRowsMatchedSQLC
	}

	return nil
}
