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

// @Summary Create a new todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param todo body services.CreateTodoRequest true "Todo details"
// @Security BearerAuth
// @Success 201 {object} TodoResponse
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos [post]
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

// @Summary List all todos
// @Tags Todo
// @Produce json
// @Security BearerAuth
// @Success 200 {array} TodoResponse
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos [get]
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

// @Summary Search todos by keyword
// @Tags Todo
// @Produce json
// @Param keyword query string true "Search keyword"
// @Security BearerAuth
// @Success 200 {array} TodoResponse
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos/search [get]
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

// @Summary Update a todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body services.UpdateTodoRequest true "Updated todo details"
// @Security BearerAuth
// @Success 200 {object} TodoResponse
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos/{id} [put]
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

// @Summary Update a todo's position
// @Tags Todo
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param position body services.UpdateTodoPositionRequest true "Updated position"
// @Security BearerAuth
// @Success 200 {object} TodoResponse
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos/{id}/position [put]
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

// @Summary Delete a todo
// @Tags Todo
// @Produce json
// @Param id path int true "Todo ID"
// @Security BearerAuth
// @Success 200 {object} gin.H "{"message": "Todo deleted"}"
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /todos/{id} [delete]
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
