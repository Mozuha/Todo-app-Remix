package router

import (
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	"todo-app/internal/middlewares"
	"todo-app/internal/services"

	"github.com/gin-gonic/gin"
)

func InitAuthHandler(sqlClient *db.Queries, passHasher services.IPasswordHasher, jwter services.ITokenGenerator) *handlers.AuthHandler {
	wrappedSqlClient := db.NewWrappedQuerier(sqlClient)
	s := services.NewAuthService(wrappedSqlClient, passHasher, jwter)
	return handlers.NewAuthHandler(s)
}

func InitUserHandler(sqlClient *db.Queries) *handlers.UserHandler {
	wrappedSqlClient := db.NewWrappedQuerier(sqlClient)
	s := services.NewUserService(wrappedSqlClient)
	return handlers.NewUserHandler(s)
}

func InitAuthMiddleware(jwter services.ITokenGenerator) gin.HandlerFunc {
	return middlewares.AuthMiddleware(jwter)
}
