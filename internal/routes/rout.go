package routes

import (
	"github.com/gofiber/fiber/v2"
	"playoo/internal/controller/rest/handlers"
)

type RoutConfig struct {
	App           *fiber.App
	AuthCheck	   fiber.Handler
	UserHandler   *handlers.UserHandler
	AuthHandler   *handlers.AuthHandler
	GameHandler   *handlers.GameHandler
	EventHandler  *handlers.EventHandler
	NewsHandler   *handlers.NewsHandler
	NoticeHandler *handlers.NoticeHandler
}

func (cfg *RoutConfig) Setup() {
	cfg.SetupCSRF()
	cfg.SetupUserRoute()
	cfg.SetupAuthRoute()
	cfg.SetupGameRoute()
	cfg.SetupNewsRoute()
	cfg.SetupEventRoute()
	cfg.SetupNotificationsRoute()
}

func (cfg *RoutConfig) SetupCSRF(){
	cfg.App.Get("/csrf",func (c *fiber.Ctx) error {
		csrf:=c.Cookies("csrf")
		return c.JSON(fiber.Map{
			"csrf_token": csrf,
		})
	})
}

func (cfg *RoutConfig) SetupUserRoute() {
	userGroup:=cfg.App.Group("/api/users", cfg.AuthCheck)

	userGroup.Patch("/avatar/:id", cfg.UserHandler.UploadAvatar)
	userGroup.Patch("/discord", cfg.UserHandler.RecordDiscord)
	userGroup.Patch("/commens", cfg.UserHandler.AddComment)
	userGroup.Patch("/follow", cfg.UserHandler.Follow)
	userGroup.Patch("/unfollow", cfg.UserHandler.Unfollow)
	userGroup.Patch("/rating", cfg.UserHandler.EditRating)

	userGroup.Delete("/avatar/:id", cfg.UserHandler.DeleteAvatar)

	userGroup.Get("/comments", cfg.UserHandler.GetComments)
	userGroup.Get("/:id", cfg.UserHandler.GetUserById)
	userGroup.Get("", cfg.UserHandler.GetUsersByAmount)
}

func (cfg *RoutConfig) SetupAuthRoute() {
	cfg.App.Post("/api/register", cfg.AuthHandler.Register)
	cfg.App.Post("/api/login", cfg.AuthHandler.Login)
	cfg.App.Use("/api/logout",cfg.AuthCheck)
	cfg.App.Post("/api/logout", cfg.AuthHandler.Logout)

	cfg.App.Get("/api/profile", cfg.AuthHandler.GetLoggedUser)
}

func (cfg *RoutConfig) SetupGameRoute() {
	gameGroup:=cfg.App.Group("/api/games", cfg.AuthCheck)

	gameGroup.Patch("",cfg.GameHandler.AddGameToUser)

	gameGroup.Delete("", cfg.GameHandler.DeleteGame)

	gameGroup.Get("", cfg.GameHandler.GetGameByName)
}

func (cfg *RoutConfig) SetupEventRoute() {
	eventsGroup:=cfg.App.Group("/api/events", cfg.AuthCheck)
	eventsGroup.Post("", cfg.EventHandler.CreateEvent)

	eventsGroup.Patch("/comments", cfg.EventHandler.AddComment)
	eventsGroup.Patch("/join", cfg.EventHandler.Join)
	eventsGroup.Patch("/unjoin", cfg.EventHandler.Unjoin)

	eventsGroup.Get("/comments", cfg.EventHandler.GetComments)
	eventsGroup.Get("/:id", cfg.EventHandler.GetEventById)

	eventsGroup.Get("/api/events", cfg.EventHandler.GetEventByAmount)
}

func (cfg *RoutConfig) SetupNewsRoute() {
	newsGroup:=cfg.App.Group("/api/news", cfg.AuthCheck)
	
	cfg.App.Post("/api/news", cfg.NewsHandler.CreateNews)

	newsGroup.Patch("/comments", cfg.NewsHandler.AddComment)

	newsGroup.Get("/comments", cfg.NewsHandler.GetComments)
	newsGroup.Get("/:id", cfg.NewsHandler.GetNewsById)

	newsGroup.Get("/api/news", cfg.NewsHandler.GetNewsByAmount)
}

func (cfg *RoutConfig) SetupNotificationsRoute() {
	notificationsGroup:=cfg.App.Group("/api/notifications", cfg.AuthCheck)

	notificationsGroup.Delete("", cfg.NoticeHandler.DeleteNotice)
	notificationsGroup.Delete("/:id", cfg.NoticeHandler.DeleteAllNotifications)

	notificationsGroup.Get("", cfg.NoticeHandler.GetNotifications)
}
