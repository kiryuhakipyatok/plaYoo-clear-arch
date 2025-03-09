package handlers

import (
	"playoo/internal/domain/service"
	"playoo/internal/dto"
	e "playoo/pkg/errors"
	"strconv"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EventHandler struct{
	EventService service.EventService
	Validator 	*validator.Validate
	Logger 		*logrus.Logger
}

func NewEventHandler(eventService service.EventService,validator *validator.Validate,logger *logrus.Logger) EventHandler{
	return EventHandler{
		EventService: eventService,
		Validator: validator,
		Logger: logger,
	}
}

func (eh EventHandler) CreateEvent(c *fiber.Ctx) error{
	ctx:=c.Context()
	request:=dto.CreateEventRequest{}
	
	if err:=c.BodyParser(&request);err!=nil{
		return e.ErrorParse(c,eh.Logger,"request",err)
	}
	if err:=eh.Validator.Struct(request);err!=nil{
		return e.FailedToValidate(c,eh.Logger,err)
	}
	event,err:=eh.EventService.CreateEvent(ctx,request.AuthorId,request.Body,request.Game,request.Max,request.Minute)
	if err!=nil{
		return e.FailedToCreate(c,eh.Logger,"event",err)
	}
	eh.Logger.Infof("event created: %v",event)
	return c.JSON(event) 
}

func (eh EventHandler) GetEventById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")
	event,err:=eh.EventService.GetById(ctx,id)
	if err!=nil{
		return e.NotFound(c,eh.Logger,"event",err)
	}
	eh.Logger.Infof("event %v received",event)
	return c.JSON(event)
}

func (eh EventHandler) GetEventByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Params("amount")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return e.ErrorParse(c,eh.Logger,"amount",err)
	}
	events,err:=eh.EventService.GetByAmount(ctx,amount)
	if err!=nil{
		return e.NotFound(c,eh.Logger,"events",err)
	}
	eh.Logger.Infof("events %v received",events)
	return c.JSON(events)
}
