package config

import (
	"playoo/internal/bot"
	"playoo/internal/controller/rest/handlers"
	"playoo/internal/domain/repository"
	"playoo/internal/domain/service"
	"playoo/internal/routes"
	"playoo/internal/shedulers"
	e "playoo/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App      		*fiber.App
	Postgres		*gorm.DB
	Redis 			*redis.Client	
	Logger 			*logrus.Logger
	Validator 		*validator.Validate
	Bot 			*bot.Bot
	ErrorHandler	*e.ErrorHandler
}

func Bootstrap(config *BootstrapConfig, stop chan struct{}){

	userRepository:=repository.NewUserRepository(config.Postgres,config.Redis)
	gameRepository:=repository.NewGameRepository(config.Postgres)
	eventRepository:=repository.NewEventRepository(config.Postgres,config.Redis)
	newsRepository:=repository.NewNewsRepository(config.Postgres)
	commentRepository:=repository.NewCommentRepository(config.Postgres)
	noticeRepository:=repository.NewNoticeRepository(config.Postgres)

	userService:=service.NewUserService(userRepository)
	authService:=service.NewAuthService(userRepository)
	gameService:=service.NewGameService(gameRepository,userRepository)
	eventService:=service.NewEventService(eventRepository,userRepository,gameRepository)
	newsService:=service.NewNewsService(newsRepository)
	noticeService:=service.NewNoticeService(noticeRepository,eventRepository,userRepository)
	commentService:=service.NewCommentService(commentRepository,userRepository,eventRepository,newsRepository)

	userHandler:=handlers.NewUserHandler(userService,commentService,config.Validator,config.Logger,config.ErrorHandler)
	authHander:=handlers.NewAuthHandler(authService,config.Validator,config.Logger,config.ErrorHandler)
	gameHandler:=handlers.NewGameHandler(gameService,config.Validator,config.Logger,config.ErrorHandler)
	eventHandler:=handlers.NewEventHandler(eventService,commentService,config.Validator,config.Logger,config.ErrorHandler)
	newsHandler:=handlers.NewNewsHandler(newsService,commentService,config.Validator,config.Logger,config.ErrorHandler)
	noticeHandler:=handlers.NewNoticeHandler(noticeService,config.Validator,config.Logger,config.ErrorHandler)

	routConfig:=routes.RoutConfig{
		App: config.App,
		UserHandler: &userHandler,
		AuthHandler: &authHander,
		GameHandler: &gameHandler,
		EventHandler: &eventHandler,
		NewsHandler: &newsHandler,
		NoticeHandler: &noticeHandler,
	}

	routConfig.Setup()
	
}

func StartShedule(config *BootstrapConfig, stop chan struct{}){
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
		Bot: config.Bot,
	}
	sheduleEvents.SetupSheduleEvents(stop)

}
