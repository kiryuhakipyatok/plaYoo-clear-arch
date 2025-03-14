package handlers

import (
	"fmt"
	"playoo/internal/domain/service"
	"playoo/internal/dto"
	e "playoo/pkg/errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EventHandler struct {
	EventService   service.EventService
	CommentService service.CommentService
	Validator      *validator.Validate
	ErrorHandler   *e.ErrorHandler
	Logger         *logrus.Logger
}

func NewEventHandler(eventService service.EventService, commentService service.CommentService, validator *validator.Validate, logger *logrus.Logger, eh *e.ErrorHandler) EventHandler {
	return EventHandler{
		EventService:   eventService,
		CommentService: commentService,
		Validator:      validator,
		Logger:         logger,
		ErrorHandler:   eh,
	}
}

func (eh *EventHandler) CreateEvent(c *fiber.Ctx) error {
	ctx := c.Context()
	request := dto.CreateEventRequest{}

	if err := c.BodyParser(&request); err != nil {
		return eh.ErrorHandler.ErrorParse(c, "request", err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return eh.ErrorHandler.FailedToValidate(c, err)
	}
	fmt.Println(request)
	event, err := eh.EventService.CreateEvent(ctx, request.AuthorId, request.Body, request.Game, request.Max, request.Minute)
	fmt.Println(err)
	if err != nil {
		return eh.ErrorHandler.FailedToCreate(c, "event", err)
	}
	eh.Logger.Infof("event created: %v", event.Id)
	response := dto.EventResponse{
		Id:       event.Id,
		AuthorId: event.AuthorId,
		Body:     event.Body,
		Game:     event.Game,
		Time:     event.Time,
		Max:      event.Max,
	}
	return c.JSON(response)
}

func (eh *EventHandler) GetEventById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	event, err := eh.EventService.GetById(ctx, id)
	if err != nil {
		return eh.ErrorHandler.NotFound(c, "event", err)
	}
	eh.Logger.Infof("event %s received", event.Id)
	return c.JSON(event)
}

func (eh *EventHandler) GetEventByAmount(c *fiber.Ctx) error {
	ctx := c.Context()
	a := c.Query("amount")
	amount, err := strconv.Atoi(a)
	if err != nil {
		return eh.ErrorHandler.ErrorParse(c, "amount", err)
	}
	events, err := eh.EventService.GetByAmount(ctx, amount)
	if err != nil {
		return eh.ErrorHandler.NotFound(c, "events", err)
	}
	eh.Logger.Infof("events %v received", events)
	return c.JSON(events)
}

func (eh *EventHandler) AddComment(c *fiber.Ctx) error {
	ctx := c.Context()
	request := dto.EventCommentRequest{}

	if err := c.BodyParser(&request); err != nil {
		return eh.ErrorHandler.ErrorParse(c, "request", err)
	}
	if err := eh.Validator.Struct(request); err != nil {
		return eh.ErrorHandler.FailedToValidate(c, err)
	}
	comment, err := eh.CommentService.AddCommentToEvent(ctx, request.AuthorId, request.ReceiverId, request.Body)
	if err != nil {
		return eh.ErrorHandler.FailedToCreate(c, "comment to event", err)
	}
	eh.Logger.Infof("comment added to event %s by %s: %s", request.ReceiverId, request.AuthorId, request.Body)
	return c.JSON(comment)
}

func (eh *EventHandler) GetComments(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Query("id")
	a := c.Query("amount")
	amount, err := strconv.Atoi(a)
	if err != nil {
		return eh.ErrorHandler.ErrorParse(c, "amount", err)
	}
	comments, err := eh.CommentService.GetComments(ctx, id, amount)
	if err != nil {
		return eh.ErrorHandler.ErrorFetching(c, "event's comments", err)
	}
	eh.Logger.Infof("comments %v received", comments)
	return c.JSON(comments)
}

func (eh *EventHandler) Join(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Query("id")
	eid := c.Query("event")
	if err := eh.EventService.Join(ctx, id, eid); err != nil {
		eh.Logger.WithError(err).Error("failed to join to event")
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to join to event",
		})
	}
	eh.Logger.Infof("user %s joined to event %s", id, eid)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (eh *EventHandler) Unjoin(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Query("id")
	eid := c.Query("event")
	if err := eh.EventService.Unjoin(ctx, id, eid); err != nil {
		eh.Logger.WithError(err).Error("failed to unjoin from event")
		c.JSON(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "failed to unjoin from event",
		})
	}
	eh.Logger.Infof("user %s unjoined from event %s", id, eid)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
