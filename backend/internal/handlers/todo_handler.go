package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	TodoService services.ITodoService
}

// Hide private userId (users.id)
type TodoResponse struct {
	ID          int32     `json:"id"`
	Description string    `json:"description"`
	Position    int64     `json:"position"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewTodoHandler(todoService services.ITodoService) *TodoHandler {
	return &TodoHandler{TodoService: todoService}
}

func (h *TodoHandler) CreateTodo(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		}
		return
	}

	var req services.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todo, err := h.TodoService.CreateTodo(ctx, userIDUuid, req)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Position:    todo.Position.Int.Int64(),
		Completed:   todo.Completed.Bool,
		CreatedAt:   todo.CreatedAt.Time,
		UpdatedAt:   todo.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusCreated, todoResponse)
}

func (h *TodoHandler) ListTodos(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		}
		return
	}

	todos, err := h.TodoService.ListTodos(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}

	todoResponses := make([]TodoResponse, len(*todos))
	for i, todo := range *todos {
		todoResponses[i] = TodoResponse{
			ID:          todo.ID,
			Description: todo.Description,
			Position:    todo.Position.Int.Int64(),
			Completed:   todo.Completed.Bool,
			CreatedAt:   todo.CreatedAt.Time,
			UpdatedAt:   todo.UpdatedAt.Time,
		}
	}

	ctx.JSON(http.StatusOK, todoResponses)
}

func (h *TodoHandler) SearchTodos(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search todos"})
		}
		return
	}

	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todos, err := h.TodoService.SearchTodos(ctx, userIDUuid, keyword)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search todos"})
		return
	}

	todoResponses := make([]TodoResponse, len(*todos))
	for i, todo := range *todos {
		todoResponses[i] = TodoResponse{
			ID:          todo.ID,
			Description: todo.Description,
			Position:    todo.Position.Int.Int64(),
			Completed:   todo.Completed.Bool,
			CreatedAt:   todo.CreatedAt.Time,
			UpdatedAt:   todo.UpdatedAt.Time,
		}
	}

	ctx.JSON(http.StatusOK, todoResponses)
}

func (h *TodoHandler) UpdateTodo(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		}
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	var req services.UpdateTodoRequest
	if reqErr := ctx.ShouldBindJSON(&req); reqErr != nil || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todo, err := h.TodoService.UpdateTodo(ctx, userIDUuid, int32(todoID), req)
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified todo not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Position:    todo.Position.Int.Int64(),
		Completed:   todo.Completed.Bool,
		CreatedAt:   todo.CreatedAt.Time,
		UpdatedAt:   todo.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, todoResponse)
}

func (h *TodoHandler) UpdateTodoPosition(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo position"})
		}
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	var req services.UpdateTodoPositionRequest
	if reqErr := ctx.ShouldBindJSON(&req); reqErr != nil || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	todo, err := h.TodoService.UpdateTodoPosition(ctx, userIDUuid, int32(todoID), req)
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified todo not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo position"})
		return
	}

	todoResponse := TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Position:    todo.Position.Int.Int64(),
		Completed:   todo.Completed.Bool,
		CreatedAt:   todo.CreatedAt.Time,
		UpdatedAt:   todo.UpdatedAt.Time,
	}

	ctx.JSON(http.StatusOK, todoResponse)
}

func (h *TodoHandler) DeleteTodo(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		}
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = h.TodoService.DeleteTodo(ctx, userIDUuid, int32(todoID))
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified todo not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
