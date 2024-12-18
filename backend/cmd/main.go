package main

import (
	"log"
	"todo-app/internal/db"
	"todo-app/internal/router"
	"todo-app/utils"
)

func main() {
	runningEnv, err := utils.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	dbpool, err := db.ConnectDB(runningEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	sqlClient := db.New(dbpool)

	redisStore, err := db.SetupRedisStore(runningEnv)
	if err != nil {
		log.Fatal(err)
	}

	r := router.SetupRouter(sqlClient, redisStore)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
