package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	middlewares_mock "todo-app/internal/middlewares/_mock"
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"
	"todo-app/internal/utils"
	"todo-app/internal/utils/testutils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/mock/gomock"
)

type testSetup struct {
	ctrl            *gomock.Controller
	mockAuthService *mock_services.MockIAuthService
	authHandler     *handlers.AuthHandler
	router          *gin.Engine
	recorder        *httptest.ResponseRecorder
	context         *gin.Context
}

func setupAuthTest(t *testing.T, useMockSession bool) *testSetup {
	ctrl := gomock.NewController(t)
	mockAuthService := mock_services.NewMockIAuthService(ctrl)
	authHandler := handlers.NewAuthHandler(mockAuthService)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	if useMockSession {
		mockSession := middlewares_mock.NewMockSession(errors.New("session save failed"))

		// Use middleware to inject the mock session
		r.Use(func(c *gin.Context) {
			// Override the default session with our mock
			c.Set(sessions.DefaultKey, mockSession)
			c.Next()
		})
	} else {
		store := memstore.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("mysession", store))
	}

	return &testSetup{
		ctrl:            ctrl,
		mockAuthService: mockAuthService,
		authHandler:     authHandler,
		router:          r,
		recorder:        w,
		context:         ctx,
	}
}

func checkSessionHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No active session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}

type want struct {
	status   int
	respFile string
}

var checkSessionWants = struct {
	exist    want
	notExist want
}{
	exist: want{
		status:   http.StatusOK,
		respFile: `{"userID": "user-id-123"}`,
	},
	notExist: want{
		status:   http.StatusUnauthorized,
		respFile: `{"error":"No active session"}`,
	},
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		reqFile        string
		want           want
		useMockSession bool
	}{
		{
			name:    "successful registration",
			reqFile: "testdata/register/201_req.json.golden",
			want: want{
				status:   http.StatusCreated,
				respFile: "testdata/register/201_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/register/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/register/400_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "user already registered",
			reqFile: "testdata/register/409_req.json.golden",
			want: want{
				status:   http.StatusConflict,
				respFile: "testdata/register/409_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/register/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/register/500_resp.json.golden",
			},
			useMockSession: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupAuthTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			// Register service won't be called when request body is invalid
			if tt.name != "invalid request body" {
				setup.mockAuthService.EXPECT().Register(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req services.RegisterRequest) (*db.User, error) {
					switch tt.want.status {
					case http.StatusCreated:
						return &db.User{}, nil
					case http.StatusConflict:
						return nil, &pgconn.PgError{Code: "23505"}
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.POST("/register", setup.authHandler.Register)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name             string
		reqFile          string
		want             want
		checkSessionWant want
		useMockSession   bool
	}{
		{
			name:    "successful login",
			reqFile: "testdata/login/200_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/login/200_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.exist,
			useMockSession:   false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/login/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/login/400_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   false,
		},
		{
			name:    "invalid email or password",
			reqFile: "testdata/login/401_req.json.golden",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/login/401_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   false,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/login/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/login/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   false,
		},
		{
			name:    "failed to save session",
			reqFile: "testdata/login/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/login/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupAuthTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			// Login service won't be called when request body is invalid
			if tt.name != "invalid request body" {
				setup.mockAuthService.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req services.LoginRequest, sessionID string) (string, string, error) {
					switch tt.want.status {
					case http.StatusOK:
						return "user-id-123", "access-token-123", nil
					case http.StatusUnauthorized:
						return "", "", utils.ErrInvalidEmailOrPswd
					case http.StatusInternalServerError:
						if tt.name == "failed to save session" {
							return "user-id-123", "access-token-123", nil
						}
						return "", "", errors.New("unexpected error")
					}
					return "", "", errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.POST("/login", setup.authHandler.Login)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))

			// Verify session
			// Simulate the check-session request with the session cookie
			cookies := setup.recorder.Result().Cookies()
			setup.recorder = httptest.NewRecorder()
			setup.context.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
			for _, cookie := range cookies {
				setup.context.Request.AddCookie(cookie) // Pass all cookies from the login response
			}

			setup.router.GET("/check-session", checkSessionHandler)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.checkSessionWant.status, []byte(tt.checkSessionWant.respFile))
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name             string
		want             want
		checkSessionWant want
		useMockSession   bool
	}{
		{
			name: "successful logout",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/logout/200_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   false,
		},
		{
			name: "failed to save session",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/logout/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.notExist,
			useMockSession:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupAuthTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			setup.context.Request = httptest.NewRequest(http.MethodPost, "/logout", nil)
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.POST("/logout", setup.authHandler.Logout)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))

			// Verify session
			cookies := setup.recorder.Result().Cookies()
			setup.recorder = httptest.NewRecorder()
			setup.context.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
			for _, cookie := range cookies {
				setup.context.Request.AddCookie(cookie)
			}

			setup.router.GET("/check-session", checkSessionHandler)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.checkSessionWant.status, []byte(tt.checkSessionWant.respFile))
		})
	}
}
