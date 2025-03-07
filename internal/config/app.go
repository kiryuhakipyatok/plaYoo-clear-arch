package config

import (
	"test/internal/bot"
	"test/internal/controller/rest/handlers"
	"test/internal/domain/repository"
	"test/internal/domain/service"
	"test/internal/routes"
	"test/internal/shedulers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App      	*fiber.App
	Postgres	*gorm.DB
	Redis 		*redis.Client	
	Bot 		*tgbotapi.BotAPI
}

func Bootstrap(config *BootstrapConfig, stop chan struct{}){

	userRepository:=repository.NewUserRepository(config.Postgres,config.Redis)
	userService:=service.NewUserService(userRepository)
	userHandler:=handlers.NewUserHandler(userService)

	authService:=service.NewAuthService(userRepository)
	authHander:=handlers.NewAuthHandler(authService)

	gameRepository:=repository.NewGameRepository(config.Postgres)
	gameService:=service.NewGameService(gameRepository,userRepository)
	gameHandler:=handlers.NewGameHandler(gameService)

	eventRepository:=repository.NewEventRepository(config.Postgres,config.Redis)
	eventService:=service.NewEventService(eventRepository,userRepository,gameRepository)
	eventHandler:=handlers.NewEventHandler(eventService)

	newsRepository:=repository.NewNewsRepository(config.Postgres)
	newsService:=service.NewNewsService(newsRepository)
	newsHandler:=handlers.NewNewsHandler(newsService)

	commentRepository:=repository.NewCommentRepository(config.Postgres)
	commentService:=service.NewCommentService(commentRepository,userRepository,eventRepository,newsRepository)
	commentHandler:=handlers.NewCommentHandler(commentService)

	noticeRepository:=repository.NewNoticeRepository(config.Postgres)
	noticeService:=service.NewNoticeService(noticeRepository,eventRepository,userRepository)
	noticeHandler:=handlers.NewNoticeHandler(noticeService)


	// bot:=bot.CreateBot(stop,userRepository)

	// sheduleEvents:=shedulers.SheduleEvents{
	// 	NoticeService: noticeService,
	// 	EventService: eventService,
	// 	UserService: userService,
	// 	Bot: bot,
	// }

	routConfig:=routes.RoutConfig{
		App: config.App,
		UserHandler: &userHandler,
		AuthHandler: &authHander,
		GameHandler: &gameHandler,
		EventHandler: &eventHandler,
		NewsHandler: &newsHandler,
		CommentHandler: &commentHandler,
		NoticeHandler: &noticeHandler,
	}

	routConfig.Setup()

	//go sheduleEvents.SetupSheduleEvents(stop)
	
}

func StartBot(config *BootstrapConfig, stop chan struct{}) *bot.Bot{
	userRepository:=repository.NewUserRepository(config.Postgres,config.Redis)
	bot:=bot.CreateBot(stop,userRepository)
	return bot
}

func StartShedule(config *BootstrapConfig, stop chan struct{}, bot *bot.Bot){
	userRepository:=repository.NewUserRepository(config.Postgres,config.Redis)
	userService:=service.NewUserService(userRepository)
	eventRepository:=repository.NewEventRepository(config.Postgres,config.Redis)
	gameRepository:=repository.NewGameRepository(config.Postgres)
	eventService:=service.NewEventService(eventRepository,userRepository,gameRepository)
	
	noticeRepository:=repository.NewNoticeRepository(config.Postgres)
	noticeService:=service.NewNoticeService(noticeRepository,eventRepository,userRepository)
	
	sheduleEvents:=shedulers.SheduleEvents{
		NoticeService: noticeService,
		EventService: eventService,
		UserService: userService,
		Bot: bot,
	}
	sheduleEvents.SetupSheduleEvents(stop)

}
