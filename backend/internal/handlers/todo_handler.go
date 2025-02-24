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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	var req services.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	todo, err := h.TodoService.CreateTodo(ctx, userIDUuid, req)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	todos, err := h.TodoService.ListTodos(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	todos, err := h.TodoService.SearchTodos(ctx, userIDUuid, keyword)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	var req services.UpdateTodoRequest
	if reqErr := ctx.ShouldBindJSON(&req); reqErr != nil || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	todo, err := h.TodoService.UpdateTodo(ctx, userIDUuid, int32(todoID), req)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	var req services.UpdateTodoPositionRequest
	if reqErr := ctx.ShouldBindJSON(&req); reqErr != nil || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	todo, err := h.TodoService.UpdateTodoPosition(ctx, userIDUuid, int32(todoID), req)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
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
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	err = h.TodoService.DeleteTodo(ctx, userIDUuid, int32(todoID))
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
