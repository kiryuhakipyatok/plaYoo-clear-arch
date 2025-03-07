package handlers

import (
	"strconv"
	"test/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler{
	return UserHandler{
		UserService: userService,
	}
}

func (uh UserHandler) GetUserById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	user,err:=uh.UserService.GetById(ctx,id)
	if err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"user not found",
		})
	}
	return c.JSON(user)
}

func (uh UserHandler) GetUsersByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")	
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error parse amount",
		})
	}
	users,err:=uh.UserService.GetByAmount(ctx,amount)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error fetch users",
		})
	}
	return c.JSON(users)
}
