package middlewares

import (
	"log"
	"net/http"
	"todo-app/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const BEARER_SCHEMA = "Bearer "

func AuthMiddleware(jwter services.TokenGenerator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bearer token will be shown like `Authorization: Bearer <token>` in http header

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			ctx.Abort()
			return
		}

		token := authHeader[len(BEARER_SCHEMA):]
		claims, err := jwter.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}
		log.Println(token)
		log.Println(claims)

		// Even if the token is valid, the user cannot be authenticated if there is no associated session
		// This prevents the user from being authenticated with valid token after logged out
		session := sessions.Default(ctx)
		sessionID := session.ID()
		userID := session.Get("userID")
		log.Println(claims.SessionID, sessionID)
		log.Println(claims.UserID, userID)
		if claims.SessionID != sessionID || claims.UserID != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
