package routes

import (
	"github.com/gofiber/fiber/v2"
	"playoo/internal/controller/rest/handlers"
)

type RoutConfig struct {
	App           *fiber.App
	UserHandler   *handlers.UserHandler
	AuthHandler   *handlers.AuthHandler
	GameHandler   *handlers.GameHandler
	EventHandler  *handlers.EventHandler
	NewsHandler   *handlers.NewsHandler
	NoticeHandler *handlers.NoticeHandler
}

func (cfg *RoutConfig) Setup() {
	cfg.SetupUserRoute()
	cfg.SetupAuthRoute()
	cfg.SetupGameRoute()
	cfg.SetupNewsRoute()
	cfg.SetupEventRoute()
	cfg.SetupNotificationsRoute()
}

func (cfg *RoutConfig) SetupUserRoute() {
	cfg.App.Patch("/api/users/avatar/:id", cfg.UserHandler.UploadAvatar)
	cfg.App.Patch("/api/users/discord", cfg.UserHandler.RecordDiscord)
	cfg.App.Patch("/api/users/commens", cfg.UserHandler.AddComment)
	cfg.App.Patch("/api/users/follow", cfg.UserHandler.Follow)
	cfg.App.Patch("/api/users/unfollow", cfg.UserHandler.Unfollow)
	cfg.App.Patch("/api/users/rating", cfg.UserHandler.EditRating)

	cfg.App.Delete("/api/users/avatar/:id", cfg.UserHandler.DeleteAvatar)

	cfg.App.Get("/api/users/comments", cfg.UserHandler.GetComments)
	cfg.App.Get("/api/users/:id", cfg.UserHandler.GetUserById)
	cfg.App.Get("/api/users", cfg.UserHandler.GetUsersByAmount)
}

func (cfg *RoutConfig) SetupAuthRoute() {
	cfg.App.Post("/api/register", cfg.AuthHandler.Register)
	cfg.App.Post("/api/login", cfg.AuthHandler.Login)
	cfg.App.Post("/api/logout", cfg.AuthHandler.Logout)

	cfg.App.Get("/api/profile", cfg.AuthHandler.GetLoggedUser)
}

func (cfg *RoutConfig) SetupGameRoute() {
	cfg.App.Patch("/api/games", cfg.GameHandler.AddGameToUser)

	cfg.App.Delete("/api/games", cfg.GameHandler.DeleteGame)

	cfg.App.Get("/api/games", cfg.GameHandler.GetGameByName)
}

func (cfg *RoutConfig) SetupEventRoute() {
	cfg.App.Post("/api/events", cfg.EventHandler.CreateEvent)

	cfg.App.Patch("/api/events/comments", cfg.EventHandler.AddComment)
	cfg.App.Patch("/api/events/join", cfg.EventHandler.Join)
	cfg.App.Patch("/api/events/unjoin", cfg.EventHandler.Unjoin)

	cfg.App.Get("/api/events/comments", cfg.EventHandler.GetComments)
	cfg.App.Get("/api/events/:id", cfg.EventHandler.GetEventById)
	cfg.App.Get("/api/events", cfg.EventHandler.GetEventByAmount)
}

func (cfg *RoutConfig) SetupNewsRoute() {
	cfg.App.Post("/api/news", cfg.NewsHandler.CreateNews)

	cfg.App.Patch("/api/news/comments", cfg.NewsHandler.AddComment)

	cfg.App.Get("/api/news/comments", cfg.NewsHandler.GetComments)
	cfg.App.Get("/api/news/:id", cfg.NewsHandler.GetNewsById)
	cfg.App.Get("/api/news", cfg.NewsHandler.GetNewsByAmount)
}

func (cfg *RoutConfig) SetupNotificationsRoute() {
	cfg.App.Delete("/api/notifications", cfg.NoticeHandler.DeleteNotice)
	cfg.App.Delete("/api/notifications/:id", cfg.NoticeHandler.DeleteAllNotifications)

	cfg.App.Get("/api/notifications", cfg.NoticeHandler.GetNotifications)
}
