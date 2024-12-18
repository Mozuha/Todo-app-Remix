package router

import (
	"os"
	"todo-app/internal/db"
	"todo-app/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func SetupRouter(sqlClient *db.Queries, redisStore redis.Store) *gin.Engine {
	r := gin.Default()
	authHandler := InitAuthHandler(sqlClient)
	userHandler := InitUserHandler(sqlClient)

	r.Use(sessions.Sessions("mysession", redisStore))

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/logout", middlewares.AuthMiddleware(), authHandler.Logout)

		users := v1.Group("/users", middlewares.AuthMiddleware())
		{
			users.GET("/me", userHandler.GetMe)
			users.PATCH("/username", userHandler.UpdateUsername)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// todos := v1.Group("/todos", middlewares.AuthMiddleware())
		// {
		// 	todos.POST("", todoHandler.CreateTodo)
		// 	todos.GET("", todoHandler.ListTodos)
		// 	todos.PATCH("/:id", todoHandler.UpdateTodo)
		// 	todos.DELETE("/:id", todoHandler.DeleteTodo)
		// 	todos.GET("/search", todoHandler.SearchTodos)
		// }
	}

	return r
}