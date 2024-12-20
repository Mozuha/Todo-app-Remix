// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteTodo(ctx context.Context, arg DeleteTodoParams) error
	DeleteUser(ctx context.Context, userID pgtype.UUID) error
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	GetUserByUserID(ctx context.Context, userID pgtype.UUID) (User, error)
	ListTodos(ctx context.Context, userID int32) ([]Todo, error)
	SearchTodos(ctx context.Context, arg SearchTodosParams) ([]Todo, error)
	UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error)
	UpdateTodoPosition(ctx context.Context, arg UpdateTodoPositionParams) (Todo, error)
	UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error
}

var _ Querier = (*Queries)(nil)
