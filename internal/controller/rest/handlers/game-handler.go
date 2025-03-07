package handlers

import (
	"strconv"
	"test/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type GameHandler struct{
	GameService service.GameService
}

func NewGameHandler(gameService service.GameService) GameHandler{
	return GameHandler{
		GameService: gameService,
	}
}

func (gh GameHandler) AddGameToTable(c *fiber.Ctx) error{
	ctx:=c.Context()
	name:=c.Query("name")
	if err:=gh.GameService.AddGameToDB(ctx,name);err!=nil{
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed add game to db",
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (gh GameHandler) AddGameToUser(c *fiber.Ctx) error{
	ctx:=c.Context()
	game:=c.Query("game")
	id:=c.Query("id")
	if err:=gh.GameService.AddGameToUser(ctx,game,id);err!=nil{
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to add game to user",
		})
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (gh GameHandler) GetGameByName(c *fiber.Ctx) error{
	ctx:=c.Context()
	name:=c.Query("game")
	game,err:=gh.GameService.GetByName(ctx,name)
	if err!=nil{
		c.JSON(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"game not found",
		})
	}
	return c.JSON(game)
}

func (gh GameHandler) GetGamesByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("game")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error parse amount",
		})
	}
	games,err:=gh.GameService.GetByAmount(ctx,amount)
	if err!=nil{
		c.JSON(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"failed fetch games",
		})
	}
	return c.JSON(games)
}