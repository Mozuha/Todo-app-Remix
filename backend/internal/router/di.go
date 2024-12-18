package router

import (
	"todo-app/internal/db"
	"todo-app/internal/handlers"
	"todo-app/internal/services"
)

func InitAuthHandler(sqlClient *db.Queries) *handlers.AuthHandler {
	s := services.NewAuthService(sqlClient)
	return handlers.NewAuthHandler(s)
}

func InitUserHandler(sqlClient *db.Queries) *handlers.UserHandler {
	s := services.NewUserService(sqlClient)
	return handlers.NewUserHandler(s)
}
