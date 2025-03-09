package routes

import (
	"playoo/internal/controller/rest/handlers"
	"github.com/gofiber/fiber/v2"
)

type RoutConfig struct {
	App 		 *fiber.App
	UserHandler  *handlers.UserHandler
	AuthHandler  *handlers.AuthHandler
	GameHandler  *handlers.GameHandler
	EventHandler *handlers.EventHandler
	NewsHandler  *handlers.NewsHandler
	CommentHandler *handlers.CommentHandler
	NoticeHandler *handlers.NoticeHandler
}

func (cfg *RoutConfig) Setup(){
	cfg.SetupUserRoute()
	cfg.SetupAuthRoute()
	cfg.SetupGameRoute()
	cfg.SetupNewsRoute()
	cfg.SetupCommentRoute()
	cfg.SetupEventRoute()
	cfg.SetupNotificationsRoute()
}

func (cfg *RoutConfig) SetupUserRoute(){
	cfg.App.Get("/api/users/:id",cfg.UserHandler.GetUserById)
	cfg.App.Get("/api/users/:amount",cfg.UserHandler.GetUsersByAmount)
}

func (cfg *RoutConfig) SetupAuthRoute(){
	cfg.App.Post("/api/register",cfg.AuthHandler.Register)
	cfg.App.Post("/api/login",cfg.AuthHandler.Login)
	cfg.App.Post("/api/logout",cfg.AuthHandler.Logout)

	cfg.App.Get("/api/profile",cfg.AuthHandler.GetLoggedUser)
}

func (cfg *RoutConfig) SetupGameRoute(){
	cfg.App.Post("/api/add-game",cfg.GameHandler.AddGameToUser)

	cfg.App.Get("/api/games",cfg.GameHandler.GetGameByName)
}

func (cfg *RoutConfig) SetupEventRoute(){
	cfg.App.Post("/api/add-event",cfg.EventHandler.CreateEvent)
	
	cfg.App.Get("/api/events/:id",cfg.EventHandler.GetEventById)
	cfg.App.Get("/api/events/:amount",cfg.EventHandler.GetEventByAmount)
}

func (cfg *RoutConfig) SetupNewsRoute(){
	cfg.App.Post("/admin/add-news",cfg.NewsHandler.CreateNews)
	
	cfg.App.Get("/api/news/:id",cfg.NewsHandler.GetNewsById)
	cfg.App.Get("/api/news/:amount",cfg.NewsHandler.GetNewsByAmount)
}

func (cfg *RoutConfig) SetupCommentRoute(){
	cfg.App.Post("/api/add-comment-to-user",cfg.CommentHandler.AddCommentToUser)
	cfg.App.Post("/api/add-comment-to-event",cfg.CommentHandler.AddCommentToEvent)
	cfg.App.Post("/api/add-comment-to-news",cfg.CommentHandler.AddCommentToNews)
}

func (cfg *RoutConfig) SetupNotificationsRoute(){
	cfg.App.Delete("/api/notifications",cfg.NoticeHandler.DeleteNotice)
	cfg.App.Delete("/api/notifications/:id",cfg.NoticeHandler.DeleteAllNotifications)
	
	cfg.App.Get("/api/notifications",cfg.NoticeHandler.GetNotifications)
}