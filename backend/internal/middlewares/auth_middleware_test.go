package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/internal/middlewares"
	middlewares_mock "todo-app/internal/middlewares/_mock"
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

type testSetup struct {
	ctrl         *gomock.Controller
	mockTokenGen *mock_services.MockITokenGenerator
	router       *gin.Engine
	recorder     *httptest.ResponseRecorder
}

var validUID = "user-id-123"
var validSID = "mock-session-id"

func setupMiddlewareTest(t *testing.T) *testSetup {
	ctrl := gomock.NewController(t)
	mockTokenGen := mock_services.NewMockITokenGenerator(ctrl)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	r := gin.New()

	// Set up mock session
	mockSession := middlewares_mock.NewMockSession(nil)
	mockSession.Set("userID", validUID)
	mockSession.Save()

	// Use middleware to inject the mock session
	r.Use(func(c *gin.Context) {
		c.Set(sessions.DefaultKey, mockSession)
		c.Next()
	})

	return &testSetup{
		ctrl:         ctrl,
		mockTokenGen: mockTokenGen,
		router:       r,
		recorder:     w,
	}
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		mockTokenResp  *services.JWTCustomClaims
		mockTokenErr   error
		expectedStatus int
	}{
		{
			name:           "valid token and session",
			authHeader:     "Bearer valid-token",
			mockTokenResp:  &services.JWTCustomClaims{UserID: validUID, SessionID: validSID},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "malformed token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid or expired token",
			authHeader:     "Bearer invalid-token",
			mockTokenErr:   errors.New("invalid or expired token"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token but invalid session",
			authHeader:     "Bearer valid-token",
			mockTokenResp:  &services.JWTCustomClaims{UserID: validUID, SessionID: "invalid-session-id"},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token with invalid user ID",
			authHeader:     "Bearer valid-token",
			mockTokenResp:  &services.JWTCustomClaims{UserID: "invalid-user-id", SessionID: validSID},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupMiddlewareTest(t)
			defer setup.ctrl.Finish()

			if tt.name != "missing auth header" && tt.name != "malformed token" {
				if tt.name == "invalid or expired token" {
					setup.mockTokenGen.EXPECT().ValidateToken("invalid-token").Return(nil, tt.mockTokenErr)
				} else {
					setup.mockTokenGen.EXPECT().ValidateToken("valid-token").Return(tt.mockTokenResp, nil)
				}
			}

			setup.router.Use(middlewares.AuthMiddleware(setup.mockTokenGen))
			setup.router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			setup.router.ServeHTTP(setup.recorder, req)

			if setup.recorder.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, setup.recorder.Code)
			}
		})
	}
}
