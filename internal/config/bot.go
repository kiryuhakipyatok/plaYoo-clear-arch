package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"playoo/internal/bot"
	"playoo/internal/domain/repository"
)

func StartBot(db *gorm.DB, redis *redis.Client, stop chan struct{}) *bot.Bot {
	userRepository := repository.NewUserRepository(db, redis)
	bot := bot.CreateBot(stop, userRepository)
	return bot

}
