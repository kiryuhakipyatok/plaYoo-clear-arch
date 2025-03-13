package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func ConnectToRedis() (*redis.Client, error) {
	ctx := context.Background()
	var (
		host     = os.Getenv("REDISHOST")
		port     = os.Getenv("REDISPORT")
		password = os.Getenv("REDISPASSWORD")
	)
	addr := fmt.Sprintf("%s:%s", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	} else {
		log.Printf("Connect to redis successfully")
	}
	return client, nil
}
