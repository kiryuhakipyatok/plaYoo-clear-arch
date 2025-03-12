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
	GameService  	 service.GameService
	Validator 		*validator.Validate
	Logger 			*logrus.Logger
	ErrorHandler	*e.ErrorHandler
}

func NewGameHandler(gameService service.GameService,validator *validator.Validate,logger *logrus.Logger, eh *e.ErrorHandler) GameHandler{
	return GameHandler{
		GameService: gameService,
		Logger: logger,
		Validator: validator,
		ErrorHandler: eh,
	}
}

func (gh *GameHandler) AddGameToUser(c *fiber.Ctx) error{
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
	gh.Logger.Infof("game %s added to %s",game,id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (gh *GameHandler) GetGameByName(c *fiber.Ctx) error{
	ctx:=c.Context()
	name:=c.Query("game")
	game,err:=gh.GameService.GetByName(ctx,name)
	if err!=nil{
		return gh.ErrorHandler.NotFound(c,"game",err)
	}
	gh.Logger.Infof("game %s received",game.Name)
	return c.JSON(game)
}

func (gh *GameHandler) GetGamesByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return gh.ErrorHandler.ErrorParse(c,"amount",err)
	}
	games,err:=gh.GameService.GetByAmount(ctx,amount)
	if err!=nil{
		return gh.ErrorHandler.ErrorFetching(c,"games",err)
	}
	gh.Logger.Infof("games %v received",games)
	return c.JSON(games)
}

func(gh *GameHandler) DeleteGame(c *fiber.Ctx) error{
	ctx:=c.Context()
	game:=c.Query("game")
	id:=c.Query("id")
	if err:=gh.GameService.DeleteGame(ctx,id,game);err!=nil{
		gh.Logger.WithError(err).Error("failed delete game from user")
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"failed delete game from user",
		})
	}
	gh.Logger.Infof("game %s deleted from user %s",game,id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}