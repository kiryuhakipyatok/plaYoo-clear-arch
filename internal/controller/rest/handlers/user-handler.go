package handlers

import (
	"strconv"
	"playoo/internal/domain/service"
	"github.com/gofiber/fiber/v2"
	e "playoo/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserHandler struct{
	UserService service.UserService
	Validator 	*validator.Validate
	Logger 		*logrus.Logger
}

func NewUserHandler(userService service.UserService,validator *validator.Validate,logger *logrus.Logger) UserHandler{
	return UserHandler{
		UserService: userService,
		Logger: logger,
		Validator: validator,
	}
}

func (uh UserHandler) GetUserById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")
	user,err:=uh.UserService.GetById(ctx,id)
	if err!=nil{
		return e.NotFound(c,uh.Logger,"user",err)
	}
	uh.Logger.Infof("user %v recieved",user)
	return c.JSON(user)
}

func (uh UserHandler) GetUsersByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Params("amount")	
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return e.ErrorParse(c,uh.Logger,"amount",err)
	}
	users,err:=uh.UserService.GetByAmount(ctx,amount)
	if err!=nil{
		return e.ErrorFetching(c,uh.Logger,"users",err)
	}
	uh.Logger.Infof("users %v recieved",users)
	return c.JSON(users)
}
