package services_test

import (
	"context"
	"errors"
	"testing"
	"todo-app/internal/db"
	mock_db "todo-app/internal/db/_mock"
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"
	"todo-app/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueries := mock_db.NewMockWrappedQuerier(ctrl)
	mockPassHasher := mock_services.NewMockIPasswordHasher(ctrl)
	mockTokenGen := mock_services.NewMockITokenGenerator(ctrl)

	authService := services.NewAuthService(mockQueries, mockPassHasher, mockTokenGen)

	uIDStr := "00010203-0405-0607-0809-0a0b0c0d0e0f"
	uIDUuid, _ := utils.StringToUUID(uIDStr)
	sessionID := "session-id-123"
	passwordCases := map[string]map[string]string{
		"correct": {
			"plain":  "password123",
			"hashed": "hashedpassword123",
		},
		"wrong": {
			"plain":  "wrongpassword",
			"hashed": "hashedwrongpassword",
		},
	}
	token := "access-token-123"

	t.Run("Register", func(t *testing.T) {
		ctx := context.Background()
		plainPassword := passwordCases["correct"]["plain"]
		hashedPassword := passwordCases["correct"]["hashed"]
		req := services.RegisterRequest{
			Email:    "test@example.com",
			Password: plainPassword,
		}

		mockPassHasher.EXPECT().
			GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost).
			Return([]byte(hashedPassword), nil)

		mockQueries.EXPECT().
			CreateUser(ctx, db.CreateUserParams{
				Email:        req.Email,
				PasswordHash: hashedPassword,
			}).
			Return(db.User{UserID: uIDUuid, Email: req.Email, PasswordHash: hashedPassword}, nil)

		user, err := authService.Register(ctx, req)

		require.NoError(t, err)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, uIDUuid, user.UserID)
	})

	t.Run("Login", func(t *testing.T) {
		ctx := context.Background()
		plainPassword := passwordCases["correct"]["plain"]
		hashedPassword := passwordCases["correct"]["hashed"]
		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: plainPassword,
		}

		mockQueries.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(db.User{UserID: uIDUuid, PasswordHash: hashedPassword}, nil)

		mockPassHasher.EXPECT().
			CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)).
			Return(nil)

		mockTokenGen.EXPECT().
			GenerateToken(uIDStr, sessionID).
			Return(token, nil)

		userID, accessToken, err := authService.Login(ctx, req, sessionID)

		require.NoError(t, err)
		assert.Equal(t, uIDStr, userID)
		assert.Equal(t, token, accessToken)
	})

	t.Run("Login_InvalidEmail", func(t *testing.T) {
		ctx := context.Background()
		plainPassword := passwordCases["correct"]["plain"]
		req := services.LoginRequest{
			Email:    "invalid@example.com",
			Password: plainPassword,
		}

		mockQueries.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(db.User{}, errors.New("user not found"))

		userID, accessToken, err := authService.Login(ctx, req, sessionID)

		assert.Error(t, err)
		assert.Equal(t, "", userID)
		assert.Equal(t, "", accessToken)
	})

	t.Run("Login_PasswordMismatch", func(t *testing.T) {
		ctx := context.Background()
		correctHashedPassword := passwordCases["correct"]["hashed"]
		wrongPlainPassword := passwordCases["wrong"]["plain"]
		req := services.LoginRequest{
			Email:    "test@example.com",
			Password: wrongPlainPassword,
		}

		mockQueries.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(db.User{PasswordHash: correctHashedPassword}, nil)

		mockPassHasher.EXPECT().
			CompareHashAndPassword([]byte(correctHashedPassword), []byte(req.Password)).
			Return(errors.New("invalid password"))

		userID, accessToken, err := authService.Login(ctx, req, sessionID)

		assert.Error(t, err)
		assert.Equal(t, "", userID)
		assert.Equal(t, "", accessToken)
	})
}
