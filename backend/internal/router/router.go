package router

import (
	"os"
	"todo-app/internal/db"
	"todo-app/internal/services"

	_ "todo-app/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Todo app API
// @version 1.0
// @license.name Apache 2.0
// @BasePath /api/v1
// @securitydefinitions.bearerauth BearerAuth
// @in header
// @name Authorization
func SetupRouter(sqlClient *db.Queries, redisStore redis.Store) *gin.Engine {
	r := gin.Default()

	passHasher := services.NewDefaultPasswordHasher()
	jwter := services.NewJWTer()
	authHandler := InitAuthHandler(sqlClient, passHasher, jwter)
	userHandler := InitUserHandler(sqlClient)
	authMiddleware := InitAuthMiddleware(jwter)
	todoHandler := InitTodoHandler(sqlClient)

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
		v1.POST("/logout", authMiddleware, authHandler.Logout)

		users := v1.Group("/me", authMiddleware)
		{
			users.GET("/", userHandler.GetMe)
			users.PATCH("/username", userHandler.UpdateMyUsername)
			users.DELETE("/", userHandler.DeleteMe)
		}

		todos := v1.Group("/todos", authMiddleware)
		{
			todos.POST("/", todoHandler.CreateTodo)
			todos.GET("/", todoHandler.ListTodos)
			todos.GET("/search", todoHandler.SearchTodos) // /search?keyword={keyword}
			todos.PUT("/:id", todoHandler.UpdateTodo)
			todos.PATCH("/:id/position", todoHandler.UpdateTodoPosition)
			todos.DELETE("/:id", todoHandler.DeleteTodo)
		}
	}

	// http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
