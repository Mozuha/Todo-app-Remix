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

// @Summary Get current user info
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GetMeResponse
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /me [get]
func (h *UserHandler) GetMe(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	user, err := h.UserService.GetMe(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, GetMeResponse{UserID: utils.UUIDToString(user.UserID), Username: user.Username, Email: user.Email})
}

// @Summary Update current user's username
// @Tags User
// @Accept json
// @Produce json
// @Param username body services.UpdateUsernameRequest true "New username"
// @Security BearerAuth
// @Success 200 {object} gin.H "{"message": "Username updated"}"
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /me/username [put]
func (h *UserHandler) UpdateMyUsername(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	var req services.UpdateUsernameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	err = h.UserService.UpdateUsername(ctx, userIDUuid, req)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Username updated"})
}

// @Summary Delete current user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gin.H "{"message": "User deleted"}"
// @Failure 404 {object} gin.H "{"error": "Resource not found"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /me [delete]
func (h *UserHandler) DeleteMe(ctx *gin.Context) {
	userIDUuid, err := utils.GetUIDFromCtxAndCreateRespUponErr(ctx)
	if err != nil {
		return
	}

	err = h.UserService.DeleteUser(ctx, userIDUuid)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrNoRowsMatchedSQLC {
			ctx.JSON(http.StatusNotFound, gin.H{"error": utils.MsgResourceNotFound})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
