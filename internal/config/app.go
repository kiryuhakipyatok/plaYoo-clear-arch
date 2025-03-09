package config

import (
	"playoo/internal/bot"
	"playoo/internal/controller/rest/handlers"
	"playoo/internal/domain/repository"
	"playoo/internal/domain/service"
	"playoo/internal/routes"
	"playoo/internal/shedulers"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	App      	*fiber.App
	Postgres	*gorm.DB
	Redis 		*redis.Client	
	Logger 		*logrus.Logger
	Validator 	*validator.Validate
	Bot 		*bot.Bot
}

func Bootstrap(config *BootstrapConfig, stop chan struct{}){


	userRepository:=repository.NewUserRepository(config.Postgres,config.Redis)
	userService:=service.NewUserService(userRepository)
	userHandler:=handlers.NewUserHandler(userService,config.Validator,config.Logger)

	authService:=service.NewAuthService(userRepository)
	authHander:=handlers.NewAuthHandler(authService,config.Validator,config.Logger)

	gameRepository:=repository.NewGameRepository(config.Postgres)
	gameService:=service.NewGameService(gameRepository,userRepository)
	gameHandler:=handlers.NewGameHandler(gameService,config.Validator,config.Logger)

	eventRepository:=repository.NewEventRepository(config.Postgres,config.Redis)
	eventService:=service.NewEventService(eventRepository,userRepository,gameRepository)
	eventHandler:=handlers.NewEventHandler(eventService,config.Validator,config.Logger)

	newsRepository:=repository.NewNewsRepository(config.Postgres)
	newsService:=service.NewNewsService(newsRepository)
	newsHandler:=handlers.NewNewsHandler(newsService,config.Validator,config.Logger)

	commentRepository:=repository.NewCommentRepository(config.Postgres)
	commentService:=service.NewCommentService(commentRepository,userRepository,eventRepository,newsRepository)
	commentHandler:=handlers.NewCommentHandler(commentService,config.Validator,config.Logger)

	noticeRepository:=repository.NewNoticeRepository(config.Postgres)
	noticeService:=service.NewNoticeService(noticeRepository,eventRepository,userRepository)
	noticeHandler:=handlers.NewNoticeHandler(noticeService,config.Validator,config.Logger)

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
