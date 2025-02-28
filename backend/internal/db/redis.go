package db

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func SetupRedisStore(runningEnv string) (redis.Store, error) {
	log.Println("opening connection to redis...")

	var connStr string
	if runningEnv == "docker" {
		connStr = os.Getenv("REDIS_ADDR")
	} else {
		connStr = os.Getenv("REDIS_ADDR_LOCALHOST")
	}

	if connStr == "" {
		return nil, fmt.Errorf("REDIS_ADDR environment variable not set")
	}

	store, err := redis.NewStore(100, "tcp", connStr, os.Getenv("REDIS_PASSWORD"), []byte(os.Getenv("REDIS_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("connected to redis")

	tokenLifeSpanHour, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_EXP_HOUR"))
	if err != nil {
		return nil, fmt.Errorf("failed to obtain required parameter for store options: %w", err)
	}

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   3600 * tokenLifeSpanHour,
		HttpOnly: true,
		Secure:   false, // TODO: Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return store, nil
}
