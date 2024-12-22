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
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"
	"todo-app/internal/utils"
	"todo-app/internal/utils/testutils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/mock/gomock"
)

type userTestSetup struct {
	ctrl            *gomock.Controller
	mockUserService *mock_services.MockIUserService
	userHandler     *handlers.UserHandler
	router          *gin.Engine
	recorder        *httptest.ResponseRecorder
	context         *gin.Context
}

var uIDStr = "00010203-0405-0607-0809-0a0b0c0d0e0f"
var uIDUuid, _ = utils.StringToUUID(uIDStr)

func setupUserTest(t *testing.T, setUserIDInCtx bool) *userTestSetup {
	ctrl := gomock.NewController(t)
	mockUserService := mock_services.NewMockIUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	if setUserIDInCtx {
		// Set a mock userID in the context for testing
		ctx.Set("userID", uIDStr)
		r.Use(func(c *gin.Context) {
			c.Set("userID", uIDStr)
			c.Next()
		})
	}

	return &userTestSetup{
		ctrl:            ctrl,
		mockUserService: mockUserService,
		userHandler:     userHandler,
		router:          r,
		recorder:        w,
		context:         ctx,
	}
}

func TestUserHandler_GetMe(t *testing.T) {
	tests := []struct {
		name           string
		want           want
		setUserIDInCtx bool
	}{
		{
			name: "successful get user",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/get_me/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name: "failed to get userID from context",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/get_me/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name: "failed to get user",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/get_me/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// GetMe service won't be called when userID is not in context
			if tt.setUserIDInCtx {
				setup.mockUserService.EXPECT().GetMe(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID) (*db.User, error) {
					switch tt.want.status {
					case http.StatusOK:
						return &db.User{UserID: uIDUuid, Username: "testuser", Email: "test@example.com"}, nil
					case http.StatusInternalServerError:
						return nil, errors.New("user not found")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodGet, "/me", nil)
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.GET("/me", setup.userHandler.GetMe)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestUserHandler_UpdateMyUsername(t *testing.T) {
	tests := []struct {
		name           string
		reqFile        string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:    "successful update username",
			reqFile: "testdata/update_my_username/200_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/update_my_username/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "failed to get userID from context",
			reqFile: "testdata/update_my_username/401_req.json.golden",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/update_my_username/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/update_my_username/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/update_my_username/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/update_my_username/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/update_my_username/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// UpdateUsername service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request body" {
				setup.mockUserService.EXPECT().UpdateUsername(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, req services.UpdateUsernameRequest) error {
					switch tt.want.status {
					case http.StatusOK:
						return nil
					case http.StatusInternalServerError:
						return errors.New("unexpected error")
					}
					return errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPut, "/me/username", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.PUT("/me/username", setup.userHandler.UpdateMyUsername)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestUserHandler_DeleteMe(t *testing.T) {
	tests := []struct {
		name           string
		want           want
		setUserIDInCtx bool
	}{
		{
			name: "successful delete user",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/delete_me/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name: "failed to get userID from context",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/delete_me/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name: "internal server error",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/delete_me/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// DeleteUser service won't be called when userID is not in context
			if tt.setUserIDInCtx {
				setup.mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID) error {
					switch tt.want.status {
					case http.StatusOK:
						return nil
					case http.StatusInternalServerError:
						return errors.New("unexpected error")
					}
					return errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodDelete, "/me", nil)
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.DELETE("/me", setup.userHandler.DeleteMe)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}
