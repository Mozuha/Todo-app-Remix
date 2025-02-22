package handlers

import (
	"log"
	"net/http"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService services.IUserService
}

type GetMeResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}

	user, err := h.UserService.GetMe(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified user not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	ctx.JSON(http.StatusOK, GetMeResponse{UserID: utils.UUIDToString(user.UserID), Username: user.Username, Email: user.Email})
}

func (h *UserHandler) UpdateMyUsername(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		}
		return
	}

	var req services.UpdateUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = h.UserService.UpdateUsername(ctx, userIDUuid, req)
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified user not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Username updated"})
}

func (h *UserHandler) DeleteMe(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCast(ctx)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "userID not found in context" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in context"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	err = h.UserService.DeleteUser(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())

		if err.Error() == "no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Specified user not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
