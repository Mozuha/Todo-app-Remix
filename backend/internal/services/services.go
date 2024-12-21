package services

import (
	"context"
	"todo-app/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type IAuthService interface {
	Register(ctx context.Context, req RegisterRequest) (db.User, error)
	Login(ctx context.Context, req LoginRequest, sessionID string) (string, string, error)
}

type IUserService interface {
	GetMe(ctx context.Context, userID pgtype.UUID) (db.User, error)
	UpdateUsername(ctx context.Context, userID pgtype.UUID, req UpdateUsernameRequest) error
	DeleteUser(ctx context.Context, userID pgtype.UUID) error
}

type IPasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type ITokenGenerator interface {
	GenerateToken(userID, sessionID string) (string, error)
	ValidateToken(tokenString string) (*JWTCustomClaims, error)
}
