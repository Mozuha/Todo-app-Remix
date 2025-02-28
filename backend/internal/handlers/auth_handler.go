package handlers

import (
	"log"
	"net/http"
	"todo-app/internal/services"
	"todo-app/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.IAuthService
}

type LoginResponse struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func NewAuthHandler(authService services.IAuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// @Summary Register an user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credential body services.RegisterRequest true "user credential"
// @Success 200 {object} gin.H "{"message": "User registered"}"
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 409 {object} gin.H "{"error": "User already registered"}"
// @Failure 500 {object} gin.H "{"error": "The server encountered unexpected error"}"
// @Router /register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req services.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	_, err := h.AuthService.Register(ctx, req)
	if err != nil {
		log.Println(err.Error())

		// TODO: Consider more manageable error handling
		if pgErr, ok := utils.AssertPgErr(err); ok {
			if pgErr.Code == "23505" {
				ctx.JSON(http.StatusConflict, gin.H{"error": "User already registered"})
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

// @Summary Login a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credential body services.LoginRequest true "user credential"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} gin.H "{"error": "Invalid request"}"
// @Failure 401 {object} gin.H "{"error": "Invalid email or password"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": utils.MsgInvalidReq})
		return
	}

	session := sessions.Default(ctx)
	sessionID := session.ID()

	userID, accessToken, err := h.AuthService.Login(ctx, req, sessionID)
	if err != nil {
		log.Println(err.Error())

		if err == utils.ErrInvalidEmailOrPswd {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": utils.MsgInvalidEmailOrPswd})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	if session.Get("userID") != nil {
		session.Clear()
		// session.Options(sessions.Options{MaxAge: -1})
		if err = session.Save(); err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
			return
		}
	}

	session.Set("userID", userID)
	if err = session.Save(); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{UserID: userID, AccessToken: accessToken})
}

// @Summary Logout a user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} gin.H "{"message": "Logged out"}"
// @Failure 500 {object} gin.H "{"error": "Internal server error"}"
// @Router /logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": utils.MsgInternalServerErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
