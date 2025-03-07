package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CreateServer() *fiber.App{
	app:=fiber.New()
	app.Static("files", "../../files")
	app.Use(cors.New(cors.Config{        
		AllowOrigins:     "*",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",         
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization", 
        ExposeHeaders:    "Content-Length",        
		AllowCredentials: false, 
    }))   
	return app
}