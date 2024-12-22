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

func setupTest(t *testing.T, useMockSession bool) *testSetup {
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

func createRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
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
	Status   int
	RespFile string
}

var checkSessionWants = struct {
	Exist    want
	NotExist want
}{
	Exist: want{
		Status:   http.StatusOK,
		RespFile: `{"userID": "user-id-123"}`,
	},
	NotExist: want{
		Status:   http.StatusUnauthorized,
		RespFile: `{"error":"No active session"}`,
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
				Status:   http.StatusCreated,
				RespFile: "testdata/register/201_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/register/400_req.json.golden",
			want: want{
				Status:   http.StatusBadRequest,
				RespFile: "testdata/register/400_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "user already registered",
			reqFile: "testdata/register/409_req.json.golden",
			want: want{
				Status:   http.StatusConflict,
				RespFile: "testdata/register/409_resp.json.golden",
			},
			useMockSession: false,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/register/500_req.json.golden",
			want: want{
				Status:   http.StatusInternalServerError,
				RespFile: "testdata/register/500_resp.json.golden",
			},
			useMockSession: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			// AuthService won't be called when request body is invalid
			if tt.name != "invalid request body" {
				setup.mockAuthService.EXPECT().Register(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req services.RegisterRequest) (*db.User, error) {
					switch tt.want.Status {
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

			setup.context.Request = createRequest(http.MethodPost, "/register", testutils.LoadFile(t, tt.reqFile))
			setup.router.POST("/register", setup.authHandler.Register)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.Status, testutils.LoadFile(t, tt.want.RespFile))
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
				Status:   http.StatusOK,
				RespFile: "testdata/login/200_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.Exist,
			useMockSession:   false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/login/400_req.json.golden",
			want: want{
				Status:   http.StatusBadRequest,
				RespFile: "testdata/login/400_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   false,
		},
		{
			name:    "invalid email or password",
			reqFile: "testdata/login/401_req.json.golden",
			want: want{
				Status:   http.StatusUnauthorized,
				RespFile: "testdata/login/401_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   false,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/login/500_req.json.golden",
			want: want{
				Status:   http.StatusInternalServerError,
				RespFile: "testdata/login/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   false,
		},
		{
			name:    "failed to save session",
			reqFile: "testdata/login/500_req.json.golden",
			want: want{
				Status:   http.StatusInternalServerError,
				RespFile: "testdata/login/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			// LoginService won't be called when request body is invalid
			if tt.name != "invalid request body" {
				setup.mockAuthService.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req services.LoginRequest, sessionID string) (string, string, error) {
					switch tt.want.Status {
					case http.StatusOK:
						return "user-id-123", "access-token-123", nil
					case http.StatusUnauthorized:
						return "", "", errors.New("invalid email or password")
					case http.StatusInternalServerError:
						if tt.name == "failed to save session" {
							return "user-id-123", "access-token-123", nil
						}
						return "", "", errors.New("unexpected error")
					}
					return "", "", errors.New("error from mock")
				})
			}

			setup.context.Request = createRequest(http.MethodPost, "/login", testutils.LoadFile(t, tt.reqFile))
			setup.router.POST("/login", setup.authHandler.Login)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.Status, testutils.LoadFile(t, tt.want.RespFile))

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

			testutils.AssertResponse(t, setup.recorder.Result(), tt.checkSessionWant.Status, []byte(tt.checkSessionWant.RespFile))
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
				Status:   http.StatusOK,
				RespFile: "testdata/logout/200_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   false,
		},
		{
			name: "failed to save session",
			want: want{
				Status:   http.StatusInternalServerError,
				RespFile: "testdata/logout/500_resp.json.golden",
			},
			checkSessionWant: checkSessionWants.NotExist,
			useMockSession:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTest(t, tt.useMockSession)
			defer setup.ctrl.Finish()

			setup.context.Request = createRequest(http.MethodPost, "/logout", nil)
			setup.router.POST("/logout", setup.authHandler.Logout)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.Status, testutils.LoadFile(t, tt.want.RespFile))

			// Verify session
			cookies := setup.recorder.Result().Cookies()
			setup.recorder = httptest.NewRecorder()
			setup.context.Request = httptest.NewRequest(http.MethodGet, "/check-session", nil)
			for _, cookie := range cookies {
				setup.context.Request.AddCookie(cookie)
			}

			setup.router.GET("/check-session", checkSessionHandler)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.checkSessionWant.Status, []byte(tt.checkSessionWant.RespFile))
		})
	}
}
