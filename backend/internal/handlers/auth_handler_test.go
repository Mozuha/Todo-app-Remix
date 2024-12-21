package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	middlewares_mock "todo-app/internal/middlewares/_mock"
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func checkSessionHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No active session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockIAuthService(ctrl)
	authHandler := handlers.NewAuthHandler(mockAuthService)

	gin.SetMode(gin.TestMode)

	t.Run("successful registration", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockAuthService.EXPECT().Register(gomock.Any(), services.RegisterRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}).Return(db.User{}, nil)

		r.POST("/register", authHandler.Register)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"message":"User registered"}`, w.Body.String())
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"invalid-email"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		r.POST("/register", authHandler.Register)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
	})

	t.Run("user already registered", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockAuthService.EXPECT().Register(gomock.Any(), services.RegisterRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}).Return(db.User{}, &pgconn.PgError{Code: "23505"})

		r.POST("/register", authHandler.Register)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"error":"User already registered"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockAuthService.EXPECT().Register(gomock.Any(), services.RegisterRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}).Return(db.User{}, errors.New("unexpected error"))

		r.POST("/register", authHandler.Register)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Failed to register user"}`, w.Body.String())
	})
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockIAuthService(ctrl)
	authHandler := handlers.NewAuthHandler(mockAuthService)

	gin.SetMode(gin.TestMode)

	t.Run("successful login", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))

		mockAuthService.EXPECT().Login(gomock.Any(), services.LoginRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}, gomock.Any()).Return("user-id-123", "access-token-123", nil)

		r.POST("/login", authHandler.Login)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handlers.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "user-id-123", resp.UserID)
		assert.Equal(t, "access-token-123", resp.AccessToken)

		// Verify session
		// Capture the session cookie from the response
		cookies := w.Result().Cookies()

		// Simulate the check-session request with the session cookie
		w = httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
		for _, cookie := range cookies {
			ctx.Request.AddCookie(cookie) // Pass all cookies from the login response
		}

		r.GET("/check-session", checkSessionHandler)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusOK, w.Code)
		var sessionResp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &sessionResp)
		assert.Equal(t, "user-id-123", sessionResp["userID"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"invalid-email"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))

		r.POST("/login", authHandler.Login)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())

		cookies := w.Result().Cookies()
		w = httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
		for _, cookie := range cookies {
			ctx.Request.AddCookie(cookie)
		}

		r.GET("/check-session", checkSessionHandler)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"No active session"}`, w.Body.String())
	})

	t.Run("invalid email or password", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"wrongpassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))

		mockAuthService.EXPECT().Login(gomock.Any(), services.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}, gomock.Any()).Return("", "", errors.New("invalid email or password"))

		r.POST("/login", authHandler.Login)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"Invalid email or password"}`, w.Body.String())

		cookies := w.Result().Cookies()
		w = httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
		for _, cookie := range cookies {
			ctx.Request.AddCookie(cookie)
		}

		r.GET("/check-session", checkSessionHandler)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"No active session"}`, w.Body.String())
	})

	t.Run("internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))

		mockAuthService.EXPECT().Login(gomock.Any(), services.LoginRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}, gomock.Any()).Return("", "", errors.New("unexpected error"))

		r.POST("/login", authHandler.Login)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Failed to log in"}`, w.Body.String())

		cookies := w.Result().Cookies()
		w = httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
		for _, cookie := range cookies {
			ctx.Request.AddCookie(cookie)
		}

		r.GET("/check-session", checkSessionHandler)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"No active session"}`, w.Body.String())
	})

	t.Run("failed to save session", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		reqBody := `{"email":"test@example.com","password":"securepassword"}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockSession := middlewares_mock.NewMockSession(errors.New("session save failed"))

		// Use middleware to inject the mock session
		r.Use(func(c *gin.Context) {
			// Override the default session with our mock
			c.Set(sessions.DefaultKey, mockSession)
			c.Next()
		})

		mockAuthService.EXPECT().Login(gomock.Any(), services.LoginRequest{
			Email:    "test@example.com",
			Password: "securepassword",
		}, gomock.Any()).Return("user-id-123", "access-token-123", nil)

		r.POST("/login", authHandler.Login)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Failed to log in"}`, w.Body.String())
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_services.NewMockIAuthService(ctrl)
	authHandler := handlers.NewAuthHandler(mockAuthService)

	gin.SetMode(gin.TestMode)

	t.Run("successful logout", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodPost, "/logout", nil)

		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))

		r.POST("/logout", authHandler.Logout)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Logged out"}`, w.Body.String())

		cookies := w.Result().Cookies()
		w = httptest.NewRecorder()
		ctx.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
		for _, cookie := range cookies {
			ctx.Request.AddCookie(cookie)
		}

		r.GET("/check-session", checkSessionHandler)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"No active session"}`, w.Body.String())
	})

	t.Run("failed to save session", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, r := gin.CreateTestContext(w)

		ctx.Request = httptest.NewRequest(http.MethodPost, "/logout", nil)

		mockSession := middlewares_mock.NewMockSession(errors.New("session save failed"))

		r.Use(func(c *gin.Context) {
			c.Set(sessions.DefaultKey, mockSession)
			c.Next()
		})

		r.POST("/logout", authHandler.Logout)
		r.ServeHTTP(w, ctx.Request)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"Failed to log out"}`, w.Body.String())
	})

}
