package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	"todo-app/internal/services"
	mock_services "todo-app/internal/services/_mock"
	"todo-app/internal/utils/testutils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/mock/gomock"
)

type todoTestSetup struct {
	ctrl            *gomock.Controller
	mockTodoService *mock_services.MockITodoService
	todoHandler     *handlers.TodoHandler
	router          *gin.Engine
	recorder        *httptest.ResponseRecorder
	context         *gin.Context
}

var mockTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func setupTodoTest(t *testing.T, setUserIDInCtx bool) *todoTestSetup {
	ctrl := gomock.NewController(t)
	mockTodoService := mock_services.NewMockITodoService(ctrl)
	todoHandler := handlers.NewTodoHandler(mockTodoService)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)

	if setUserIDInCtx {
		ctx.Set("userID", uIDStr)
		r.Use(func(c *gin.Context) {
			c.Set("userID", uIDStr)
			c.Next()
		})
	}

	return &todoTestSetup{
		ctrl:            ctrl,
		mockTodoService: mockTodoService,
		todoHandler:     todoHandler,
		router:          r,
		recorder:        w,
		context:         ctx,
	}
}

func TestTodoHandler_CreateTodo(t *testing.T) {
	tests := []struct {
		name           string
		reqFile        string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:    "successful create todo",
			reqFile: "testdata/create_todo/201_req.json.golden",
			want: want{
				status:   http.StatusCreated,
				respFile: "testdata/create_todo/201_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "failed to get userID from context",
			reqFile: "testdata/create_todo/401_req.json.golden",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/create_todo/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:    "invalid request body",
			reqFile: "testdata/create_todo/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/create_todo/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "internal server error",
			reqFile: "testdata/create_todo/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/create_todo/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// CreateTodo service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request body" {
				setup.mockTodoService.EXPECT().CreateTodo(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, req services.CreateTodoRequest) (*db.Todo, error) {
					switch tt.want.status {
					case http.StatusCreated:
						return &db.Todo{
							ID:          1,
							Description: req.Description,
							Position:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
							Completed:   pgtype.Bool{Bool: false, Valid: true},
							CreatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
							UpdatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
						}, nil
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.POST("/todos", setup.todoHandler.CreateTodo)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestTodoHandler_ListTodos(t *testing.T) {
	tests := []struct {
		name           string
		want           want
		setUserIDInCtx bool
	}{
		{
			name: "successful list todos",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/list_todos/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name: "successful list todos - empty list",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/list_todos/200_resp_empty.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name: "failed to get userID from context",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/list_todos/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name: "internal server error",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/list_todos/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			if tt.setUserIDInCtx {
				setup.mockTodoService.EXPECT().ListTodos(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID) (*[]db.Todo, error) {
					switch tt.want.status {
					case http.StatusOK:
						if tt.name != "successful list todos - empty list" {
							return &[]db.Todo{{
								ID:          1,
								Description: "Test todo",
								Position:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
								Completed:   pgtype.Bool{Bool: false, Valid: true},
								CreatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
								UpdatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
							}}, nil
						} else {
							return &[]db.Todo{}, nil
						}
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodGet, "/todos", nil)
			setup.router.GET("/todos", setup.todoHandler.ListTodos)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestTodoHandler_SearchTodos(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:       "successful search todos",
			queryParam: "keyword=Test",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/search_todos/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:       "successful search todos - empty list",
			queryParam: "keyword=Void",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/search_todos/200_resp_empty.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:       "failed to get userID from context",
			queryParam: "keyword=Test",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/search_todos/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:       "invalid request",
			queryParam: "invalid_field=invalid_value",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/search_todos/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:       "internal server error",
			queryParam: "keyword=Test",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/search_todos/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// SearchTodos service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request" {
				setup.mockTodoService.EXPECT().SearchTodos(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, keyword string) (*[]db.Todo, error) {
					switch tt.want.status {
					case http.StatusOK:
						if tt.name != "successful search todos - empty list" {
							return &[]db.Todo{{
								ID:          1,
								Description: "Test todo",
								Position:    pgtype.Numeric{Int: big.NewInt(100), Valid: true},
								Completed:   pgtype.Bool{Bool: false, Valid: true},
								CreatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
								UpdatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
							}}, nil
						} else {
							return &[]db.Todo{}, nil
						}
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodGet, "/todos/search?"+tt.queryParam, nil)
			setup.router.GET("/todos/search", setup.todoHandler.SearchTodos)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestTodoHandler_UpdateTodo(t *testing.T) {
	tests := []struct {
		name           string
		todoID         string
		reqFile        string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:    "successful update todo",
			todoID:  "1",
			reqFile: "testdata/update_todo/200_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/update_todo/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "failed to get userID from context",
			todoID:  "1",
			reqFile: "testdata/update_todo/401_req.json.golden",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/update_todo/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:    "invalid request body",
			todoID:  "1",
			reqFile: "testdata/update_todo/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/update_todo/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "specified todo not found",
			todoID:  "1000",
			reqFile: "testdata/update_todo/404_req.json.golden",
			want: want{
				status:   http.StatusNotFound,
				respFile: "testdata/update_todo/404_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "internal server error",
			todoID:  "1",
			reqFile: "testdata/update_todo/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/update_todo/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// UpdateTodo service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request body" {
				setup.mockTodoService.EXPECT().UpdateTodo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, todoID int32, req services.UpdateTodoRequest) (*db.Todo, error) {
					switch tt.want.status {
					case http.StatusOK:
						return &db.Todo{
							ID:          todoID,
							Description: req.Description,
							Position:    pgtype.Numeric{Int: big.NewInt(req.Position), Valid: true},
							Completed:   pgtype.Bool{Bool: req.Completed, Valid: true},
							CreatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
							UpdatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
						}, nil
					case http.StatusNotFound:
						return nil, errors.New("no rows in result set")
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPut, "/todos/"+tt.todoID, bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.PUT("/todos/:id", setup.todoHandler.UpdateTodo)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

// Let's say a todo C[pos=300] has been moved inbetween A[pos=100] and B[pos=200]
func TestTodoHandler_UpdateTodoPosition(t *testing.T) {
	tests := []struct {
		name           string
		todoID         string
		reqFile        string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:    "successful update todo position",
			todoID:  "3",
			reqFile: "testdata/update_todo_position/200_req.json.golden",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/update_todo_position/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "failed to get userID from context",
			todoID:  "3",
			reqFile: "testdata/update_todo_position/401_req.json.golden",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/update_todo_position/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:    "invalid request",
			todoID:  "3",
			reqFile: "testdata/update_todo_position/400_req.json.golden",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/update_todo_position/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "specified todo not found",
			todoID:  "1000",
			reqFile: "testdata/update_todo_position/404_req.json.golden",
			want: want{
				status:   http.StatusNotFound,
				respFile: "testdata/update_todo_position/404_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:    "internal server error",
			todoID:  "3",
			reqFile: "testdata/update_todo_position/500_req.json.golden",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/update_todo_position/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// UpdateTodoPosition service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request" {
				setup.mockTodoService.EXPECT().UpdateTodoPosition(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, todoID int32, req services.UpdateTodoPositionRequest) (*db.Todo, error) {
					switch tt.want.status {
					case http.StatusOK:
						return &db.Todo{
							ID:          todoID,
							Description: "Updated todo",
							Position:    pgtype.Numeric{Int: big.NewInt(150), Valid: true},
							Completed:   pgtype.Bool{Bool: false, Valid: true},
							CreatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
							UpdatedAt:   pgtype.Timestamptz{Time: mockTime, Valid: true},
						}, nil
					case http.StatusNotFound:
						return nil, errors.New("no rows in result set")
					case http.StatusInternalServerError:
						return nil, errors.New("unexpected error")
					}
					return nil, errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodPatch, "/todos/"+tt.todoID+"/position", bytes.NewReader(testutils.LoadFile(t, tt.reqFile)))
			setup.context.Request.Header.Set("Content-Type", "application/json")
			setup.router.PATCH("/todos/:id/position", setup.todoHandler.UpdateTodoPosition)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}

func TestTodoHandler_DeleteTodo(t *testing.T) {
	tests := []struct {
		name           string
		todoID         string
		want           want
		setUserIDInCtx bool
	}{
		{
			name:   "successful delete todo",
			todoID: "1",
			want: want{
				status:   http.StatusOK,
				respFile: "testdata/delete_todo/200_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:   "failed to get userID from context",
			todoID: "1",
			want: want{
				status:   http.StatusUnauthorized,
				respFile: "testdata/delete_todo/401_resp.json.golden",
			},
			setUserIDInCtx: false,
		},
		{
			name:   "invalid request",
			todoID: "invalid",
			want: want{
				status:   http.StatusBadRequest,
				respFile: "testdata/delete_todo/400_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:   "specified todo not found",
			todoID: "1000",
			want: want{
				status:   http.StatusNotFound,
				respFile: "testdata/delete_todo/404_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
		{
			name:   "internal server error",
			todoID: "1",
			want: want{
				status:   http.StatusInternalServerError,
				respFile: "testdata/delete_todo/500_resp.json.golden",
			},
			setUserIDInCtx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := setupTodoTest(t, tt.setUserIDInCtx)
			defer setup.ctrl.Finish()

			// DeleteTodo service won't be called when userID is not in context or request body is invalid
			if tt.setUserIDInCtx && tt.name != "invalid request" {
				setup.mockTodoService.EXPECT().DeleteTodo(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, userID pgtype.UUID, todoID int32) error {
					switch tt.want.status {
					case http.StatusOK:
						return nil
					case http.StatusNotFound:
						return errors.New("no rows in result set")
					case http.StatusInternalServerError:
						return errors.New("unexpected error")
					}
					return errors.New("error from mock")
				})
			}

			setup.context.Request = httptest.NewRequest(http.MethodDelete, "/todos/"+tt.todoID, nil)
			setup.router.DELETE("/todos/:id", setup.todoHandler.DeleteTodo)
			setup.router.ServeHTTP(setup.recorder, setup.context.Request)

			testutils.AssertResponse(t, setup.recorder.Result(), tt.want.status, testutils.LoadFile(t, tt.want.respFile))
		})
	}
}
