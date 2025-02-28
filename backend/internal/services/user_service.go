package services

import (
	"context"
	"todo-app/internal/db"
	"todo-app/internal/utils"

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

func (s *UserService) GetMe(ctx context.Context, userID pgtype.UUID) (*db.User, error) {
	user, err := s.SqlClient.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUsername(ctx context.Context, userID pgtype.UUID, req UpdateUsernameRequest) error {
	return s.SqlClient.UpdateUsername(ctx, db.UpdateUsernameParams{
		Username: req.Username,
		UserID:   userID,
	})
}

func (s *UserService) DeleteUser(ctx context.Context, userID pgtype.UUID) error {
	user, err := s.SqlClient.DeleteUser(ctx, userID)
	if err != nil {
		return err
	} else if user.ID == 0 {
		return utils.ErrNoRowsMatchedSQLC
	}

	return nil
}
