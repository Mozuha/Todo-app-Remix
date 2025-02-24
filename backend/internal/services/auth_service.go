package services

import (
	"context"
	"errors"
	"todo-app/internal/db"
	"todo-app/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	SqlClient      db.WrappedQuerier
	PasswordHasher IPasswordHasher
	TokenGenerator ITokenGenerator
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthService(sqlClient db.WrappedQuerier, passHasher IPasswordHasher, jwter ITokenGenerator) *AuthService {
	return &AuthService{SqlClient: sqlClient, PasswordHasher: passHasher, TokenGenerator: jwter}
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*db.User, error) {
	hashedPassword, err := s.PasswordHasher.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.SqlClient.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest, sessionID string) (string, string, error) {
	user, err := s.SqlClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", "", utils.ErrInvalidEmailOrPswd
	}

	if err = s.PasswordHasher.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", "", utils.ErrInvalidEmailOrPswd
	}

	userIDStr := utils.UUIDToString(user.UserID)
	if userIDStr == "" {
		return "", "", errors.New("failed to convert uuid to string")
	}

	accessToken, err := s.TokenGenerator.GenerateToken(userIDStr, sessionID)
	if err != nil {
		return "", "", err
	}

	return userIDStr, accessToken, nil
}
