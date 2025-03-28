package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"playoo/internal/bot"
	"playoo/internal/controller/rest/handlers"
	"playoo/internal/domain/repository"
	"playoo/internal/domain/service"
	"playoo/internal/routes"
	"playoo/internal/shedulers"
	e "playoo/pkg/errors"
	//"github.com/gofiber/jwt/v3"
	//"os"
)

type BootstrapConfig struct {
	App          *fiber.App
	Postgres     *gorm.DB
	Redis        *redis.Client
	Logger       *logrus.Logger
	Validator    *validator.Validate
	Bot          *bot.Bot
	ErrorHandler *e.ErrorHandler
}

func Bootstrap(config *BootstrapConfig, stop chan struct{}) {

	transactor := repository.NewTransactor(config.Postgres)

	userRepository := repository.NewUserRepository(config.Postgres, config.Redis)
	gameRepository := repository.NewGameRepository(config.Postgres)
	eventRepository := repository.NewEventRepository(config.Postgres, config.Redis)
	newsRepository := repository.NewNewsRepository(config.Postgres)
	commentRepository := repository.NewCommentRepository(config.Postgres)
	noticeRepository := repository.NewNoticeRepository(config.Postgres)

	userService := service.NewUserService(userRepository, transactor)
	authService := service.NewAuthService(userRepository)
	gameService := service.NewGameService(gameRepository, userRepository, transactor)
	eventService := service.NewEventService(eventRepository, userRepository, gameRepository, transactor)
	newsService := service.NewNewsService(newsRepository)
	noticeService := service.NewNoticeService(noticeRepository, eventRepository, userRepository, transactor)
	commentService := service.NewCommentService(commentRepository, userRepository, eventRepository, newsRepository, transactor)

	userHandler := handlers.NewUserHandler(userService, commentService, config.Validator, config.Logger, config.ErrorHandler)
	authHander := handlers.NewAuthHandler(authService, config.Validator, config.Logger, config.ErrorHandler)
	gameHandler := handlers.NewGameHandler(gameService, config.Validator, config.Logger, config.ErrorHandler)
	eventHandler := handlers.NewEventHandler(eventService, commentService, config.Validator, config.Logger, config.ErrorHandler)
	newsHandler := handlers.NewNewsHandler(newsService, commentService, config.Validator, config.Logger, config.ErrorHandler)
	noticeHandler := handlers.NewNoticeHandler(noticeService, config.Validator, config.Logger, config.ErrorHandler)

	// jwtMiddleware := jwtware.New(jwtware.Config{
    //     SigningKey: []byte(os.Getenv("SECRET")),
    //     ErrorHandler: func(c *fiber.Ctx, err error) error {
	// 		c.Status(fiber.StatusUnauthorized)
    //         return c.JSON(fiber.Map{
    //             "message": "unauthorized",
    //         })
    //     },
    // })

	routConfig := routes.RoutConfig{
		App:           config.App,
		//AuthCheck: 	   jwtMiddleware,
		UserHandler:   &userHandler,
		AuthHandler:   &authHander,
		GameHandler:   &gameHandler,
		EventHandler:  &eventHandler,
		NewsHandler:   &newsHandler,
		NoticeHandler: &noticeHandler,
	}

	routConfig.Setup()

}

func StartShedule(config *BootstrapConfig, stop chan struct{}) {
	transactor := repository.NewTransactor(config.Postgres)
	userRepository := repository.NewUserRepository(config.Postgres, config.Redis)
	userService := service.NewUserService(userRepository, transactor)
	eventRepository := repository.NewEventRepository(config.Postgres, config.Redis)
	gameRepository := repository.NewGameRepository(config.Postgres)
	eventService := service.NewEventService(eventRepository, userRepository, gameRepository, transactor)

	noticeRepository := repository.NewNoticeRepository(config.Postgres)
	noticeService := service.NewNoticeService(noticeRepository, eventRepository, userRepository, transactor)

	sheduleEvents := shedulers.SheduleEvents{
		NoticeService: noticeService,
		EventService:  eventService,
		UserService:   userService,
		Logger: 	   config.Logger,
		Bot:           config.Bot,
	}
	sheduleEvents.SetupSheduleEvents(stop)

}
