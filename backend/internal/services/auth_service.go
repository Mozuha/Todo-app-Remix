package services

import (
	"context"
	"errors"
	"todo-app/internal/db"
	"todo-app/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SqlClient *db.Queries
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthService(sqlClient *db.Queries) *AuthService {
	return &AuthService{SqlClient: sqlClient}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}

	user, err := s.SqlClient.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return db.User{}, err
	}

	return user, err
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest, sessionID string) (string, string, error) {
	user, err := s.SqlClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	userIDStr := utils.UUIDToString(user.UserID)
	if userIDStr == "" {
		return "", "", errors.New("failed to convert uuid to string")
	}

	accessToken, err := utils.GenerateJWT(userIDStr, sessionID)
	if err != nil {
		return "", "", err
	}

	return userIDStr, accessToken, nil
}
