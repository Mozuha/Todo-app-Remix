package router

import (
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	"todo-app/internal/middlewares"
	"todo-app/internal/services"

	"github.com/gin-gonic/gin"
)

func InitAuthHandler(sqlClient *db.Queries, passHasher services.PasswordHasher, jwter services.TokenGenerator) *handlers.AuthHandler {
	wrappedSqlClient := db.NewWrappedQuerier(sqlClient)
	s := services.NewAuthService(wrappedSqlClient, passHasher, jwter)
	return handlers.NewAuthHandler(s)
}

func InitUserHandler(sqlClient *db.Queries) *handlers.UserHandler {
	wrappedSqlClient := db.NewWrappedQuerier(sqlClient)
	s := services.NewUserService(wrappedSqlClient)
	return handlers.NewUserHandler(s)
}

func InitAuthMiddleware(jwter services.TokenGenerator) gin.HandlerFunc {
	return middlewares.AuthMiddleware(jwter)
}
