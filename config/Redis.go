package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

var (
	RedisClient *redis.Client
)

func SetRedisStore(router *gin.Engine) {
	poolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	rdb := redis.NewClient(&redis.Options{
		Addr:             redisAddress,
		Password:         redisPassword,
		DB:               redisDB,
		DisableIndentity: true,
		PoolSize:         poolSize,
		MinIdleConns:     3,
		DialTimeout:      5 * time.Second,
		ReadTimeout:      3 * time.Second,
		WriteTimeout:     3 * time.Second,
	})
	pong, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		panic("failed to connect to Redis" + err.Error())
		return
	}
	fmt.Println("Connected to Redis:", pong)
	RedisClient = rdb
}
