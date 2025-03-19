package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/utils"
	"time"
)

func CreateServer() *fiber.App {
	app := fiber.New()
	app.Static("files", "../../files")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
	}))
	app.Use(csrf.New(csrf.Config{
        KeyLookup:      "header:X-CSRF-Token",
        CookieName:     "csrf_",
        CookieSameSite: "Lax",
        Expiration:     1 * time.Hour,
        KeyGenerator:   utils.UUID,
    }))
	return app
}
