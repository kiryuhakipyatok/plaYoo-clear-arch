package handlers

import (
	"strconv"
	"test/internal/domain/service"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct{
	EventService service.EventService
}

func NewEventHandler(eventService service.EventService) EventHandler{
	return EventHandler{
		EventService: eventService,
	}
}

func (eh EventHandler) CreateEvent(c *fiber.Ctx) error{
	ctx:=c.Context()
	var request struct{
		AuthorId 	string 	`json:"author-id"`
		Game 		string 	`json:"game"`
		Body 		string  `json:"body"`
		Max 		int 	`json:"max"`
		Minute 		int 	`json:"minute"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"falied to parse request",
		})
	}
	event,err:=eh.EventService.CreateEvent(ctx,request.AuthorId,request.Body,request.Game,request.Max,request.Minute)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"falied to create event",
		})
	}
	return c.JSON(event)
}

func (eh EventHandler) GetEventById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	event,err:=eh.EventService.GetById(ctx,id)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"event not found",
		})
	}
	return c.JSON(event)
}

func (eh EventHandler) GetEventByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error parse amount",
		})
	}
	event,err:=eh.EventService.GetByAmount(ctx,amount)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"event not found",
		})
	}
	return c.JSON(event)
}
