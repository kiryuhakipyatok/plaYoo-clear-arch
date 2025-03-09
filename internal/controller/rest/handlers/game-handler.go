package handlers

import (
	"strconv"
	"playoo/internal/domain/service"
	"github.com/gofiber/fiber/v2"
	e "playoo/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type GameHandler struct{
	GameService service.GameService
	Validator 	*validator.Validate
	Logger 		*logrus.Logger
}

func NewGameHandler(gameService service.GameService,validator *validator.Validate,logger *logrus.Logger) GameHandler{
	return GameHandler{
		GameService: gameService,
		Logger: logger,
		Validator: validator,
	}
}

func (gh GameHandler) AddGameToUser(c *fiber.Ctx) error{
	ctx:=c.Context()
	game:=c.Query("game")
	id:=c.Query("id")
	if err:=gh.GameService.AddGameToUser(ctx,game,id);err!=nil{
		gh.Logger.WithError(err).Error("failed add game to user")
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed to add game to user",
		})
	}
	gh.Logger.Infof("game %v added to %v",game,id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (gh GameHandler) GetGameByName(c *fiber.Ctx) error{
	ctx:=c.Context()
	name:=c.Query("game")
	game,err:=gh.GameService.GetByName(ctx,name)
	if err!=nil{
		return e.NotFound(c,gh.Logger,"game",err)
	}
	gh.Logger.Infof("game %v received",game)
	return c.JSON(game)
}

func (gh GameHandler) GetGamesByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("game")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return e.ErrorParse(c,gh.Logger,"amount",err)
	}
	games,err:=gh.GameService.GetByAmount(ctx,amount)
	if err!=nil{
		return e.ErrorFetching(c,gh.Logger,"games",err)
	}
	gh.Logger.Infof("games %v received",games)
	return c.JSON(games)
}