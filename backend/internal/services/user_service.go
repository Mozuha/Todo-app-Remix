package services

import (
	"context"
	"todo-app/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	SqlClient db.WrappedQuerier
}

type UpdateUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

func NewUserService(sqlClient db.WrappedQuerier) *UserService {
	return &UserService{SqlClient: sqlClient}
}

func (s *UserService) GetMe(ctx context.Context, userID pgtype.UUID) (db.User, error) {
	return s.SqlClient.GetUserByUserID(ctx, userID)
}

func (s *UserService) UpdateUsername(ctx context.Context, userID pgtype.UUID, req UpdateUsernameRequest) error {
	return s.SqlClient.UpdateUsername(ctx, db.UpdateUsernameParams{
		Username: req.Username,
		UserID:   userID,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, userID pgtype.UUID) error {
	return s.SqlClient.DeleteUser(ctx, userID)
}
