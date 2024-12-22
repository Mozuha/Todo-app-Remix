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

func setupUserTest(t *testing.T) *userTestSetup {
	ctrl := gomock.NewController(t)
	mockUserService := mock_services.NewMockIUserService(ctrl)
	userHandler := handlers.NewUserHandler(mockUserService)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	// Set a mock userID in the context for testing
	ctx.Set("userID", uIDStr)
	r.Use(func(c *gin.Context) {
		c.Set("userID", uIDStr)
		c.Next()
	})

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
		name string
		want want
	}{
		{
			name: "successful get user",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/getme/200_resp.json.golden",
			},
		},
		{
			name: "failed to get user",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/getme/500_resp.json.golden",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t)
			defer setup.ctrl.Finish()

			setup.mockUserService.EXPECT().GetMe(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID) (*db.User, error) {
				switch tt.want.status {
				case http.StatusOK:
					return &db.User{UserID: uIDUuid, Username: "testuser", Email: "test@example.com"}, nil
				case http.StatusInternalServerError:
					return nil, errors.New("user not found")
				}
				return nil, errors.New("error from mock")
			})

			setup.context.Request = httptest.NewRequest(http.MethodGet, "/users/me", nil)
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.GET("/users/me", setup.userHandler.GetMe)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestUserHandler_UpdateUsername(t *testing.T) {
	tests := []struct {
		name    string
		reqFile string
		want    want
	}{
		{
			name:    "successful update username",
			reqFile: "testdata/update_username/200_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/update_username/200_resp.json.golden",
			},
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/update_username/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/update_username/400_resp.json.golden",
			},
		},
		{
			name:    "internal server error",
			reqFile: "testdata/update_username/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/update_username/500_resp.json.golden",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t)
			defer setup.ctrl.Finish()

			// UpdateUsername service won't be called when request body is invalid
			if tt.name != "invalid request body" {
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

			setup.context.Request = httptest.NewRequest(http.MethodPut, "/users/username", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.PUT("/users/username", setup.userHandler.UpdateUsername)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name string
		want want
	}{
		{
			name: "successful delete user",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/delete_user/200_resp.json.golden",
			},
		},
		{
			name: "internal server error",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/delete_user/500_resp.json.golden",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupUserTest(t)
			defer setup.ctrl.Finish()

			setup.mockUserService.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID) error {
				switch tt.want.status {
				case http.StatusOK:
					return nil
				case http.StatusInternalServerError:
					return errors.New("unexpected error")
				}
				return errors.New("error from mock")
			})

			setup.context.Request = httptest.NewRequest(http.MethodDelete, "/users/me", nil)
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.DELETE("/users/me", setup.userHandler.DeleteUser)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}
