package main

import (
	"context"
	"fmt"
	"log"
	"todo-app/internal/db"
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
	insertedUser, err := sqlClient.CreateUser(context.Background(), db.CreateUserParams{Email: "sample@sample.com", PasswordHash: "sample"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedUser)
}
