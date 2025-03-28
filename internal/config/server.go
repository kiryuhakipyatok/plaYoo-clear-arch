package config

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v3"
	"os"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/utils"
)

func CreateServer() *fiber.App {
	app := fiber.New()
	app.Static("files", "/files")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,

	}), func (c *fiber.Ctx) error {
		if c.Path() == "/api/register" || c.Path() == "/api/login" || c.Path() == "/api/logout"{
			return c.Next()
		}
		return csrf.New(csrf.Config{
			KeyLookup:      "headers:X-CSRF-Token",
			CookieName:     "csrf",
			CookieSameSite: "Lax",
			Expiration:     1 * time.Hour,
			KeyGenerator:   utils.UUID,
			CookieHTTPOnly: true,
		})(c)
	}, func (c *fiber.Ctx) error {
		if c.Path() == "/api/register" || c.Path() == "/api/login"{
			return c.Next()
		}
		return jwtware.New(jwtware.Config{
			SigningKey: []byte(os.Getenv("SECRET")),
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				c.Status(fiber.StatusUnauthorized)
				return c.JSON(fiber.Map{
					"message": "unauthorized",
				})
			},
			TokenLookup: "cookie:jwt",
		})(c)
	})

	
	// app.Get("/csrf",func (c *fiber.Ctx) error  {
	// 	return c.JSON(fiber.Map{
	// 		"csrf":"csrf in cookie",
	// 	})
	// })
	return app
}
