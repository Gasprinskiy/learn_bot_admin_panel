package main

import (
	"context"
	"learn_bot_admin_panel/config"
	"log"
	"os"
	"os/signal"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := config.NewConfig()

	// подключение к redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})
	defer rdb.Close()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Panic("ошибка при пинге redis: ", err)
	}
}
