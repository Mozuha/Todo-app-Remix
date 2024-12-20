package handlers

import (
	"log"
	"net/http"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	UserService *services.UserService
}

type GetMeResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func getUIDFromCtxAndCast(ctx *gin.Context) pgtype.UUID {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Required data not found"})
	}

	userIDUuid, err := utils.StringToUUID(userID.(string))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert user_id to expected type"})
	}

	return userIDUuid
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	userIDUuid := getUIDFromCtxAndCast(ctx)

	user, err := h.UserService.GetMe(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, GetMeResponse{UserID: utils.UUIDToString(user.UserID), Username: user.Username, Email: user.Email})
}

func (h *UserHandler) UpdateUsername(ctx *gin.Context) {
	userIDUuid := getUIDFromCtxAndCast(ctx)

	var req services.UpdateUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.UserService.UpdateUsername(ctx, userIDUuid, req)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Username updated"})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userIDUuid := getUIDFromCtxAndCast(ctx)

	err := h.UserService.DeleteUser(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
