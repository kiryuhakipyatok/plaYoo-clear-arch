package handlers

import (
	"playoo/internal/domain/service"
	"playoo/internal/dto"
	e "playoo/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CommentHandler struct{
	CommentService 	 service.CommentService
	Validator 		*validator.Validate
	Logger 			*logrus.Logger
}

func NewCommentHandler(commentService service.CommentService,validator *validator.Validate,logger *logrus.Logger) CommentHandler{
	return CommentHandler{
		CommentService: commentService,
		Validator: validator,
		Logger: logger,
	}
}


func (ch *CommentHandler) AddCommentToUser(c *fiber.Ctx) error{
	ctx:=c.Context()
	request:=dto.UserCommentRequest{}
	
	if err:=c.BodyParser(&request);err!=nil{
		return e.ErrorParse(c,ch.Logger,"request",err)
	}
	if err:=ch.Validator.Struct(request);err!=nil{

		return e.FailedToValidate(c,ch.Logger,err)
	}
	comment,err:=ch.CommentService.AddCommentToUser(ctx,request.AuthorId,request.ReceiverId,request.Body)
	if err!=nil{
		return e.FailedToCreate(c,ch.Logger,"comment to user",err)
	}
	ch.Logger.Infof("comment added to user %s by %s: %s",request.ReceiverId,request.AuthorId,request.Body)
	return c.JSON(comment)
}

func (ch *CommentHandler) AddCommentToEvent(c *fiber.Ctx) error{
	ctx:=c.Context()
	request:=dto.EventCommentRequest{}
	
	if err:=c.BodyParser(&request);err!=nil{
		return e.ErrorParse(c,ch.Logger,"request",err)
	}
	if err:=ch.Validator.Struct(request);err!=nil{
		return e.FailedToValidate(c,ch.Logger,err)
	}
	comment,err:=ch.CommentService.AddCommentToEvent(ctx,request.AuthorId,request.ReceiverId,request.Body)
	if err!=nil{
		return e.FailedToCreate(c,ch.Logger,"comment to event",err)
	}
	ch.Logger.Infof("comment added to event %s by %s: %s",request.ReceiverId,request.AuthorId,request.Body)
	return c.JSON(comment)
}

func (ch *CommentHandler) AddCommentToNews(c *fiber.Ctx) error{
	ctx:=c.Context()
	request:=dto.NewsCommentRequest{}
	
	if err:=c.BodyParser(&request);err!=nil{
		return e.ErrorParse(c,ch.Logger,"request",err)
	}
	if err:=ch.Validator.Struct(request);err!=nil{
		return e.FailedToValidate(c,ch.Logger,err)
	}
	comment,err:=ch.CommentService.AddCommentToUser(ctx,request.AuthorId,request.ReceiverId,request.Body)
	if err!=nil{
		return e.FailedToCreate(c,ch.Logger,"comment to news",err)
	}
	ch.Logger.Infof("comment added to news %s by %s: %s",request.ReceiverId,request.AuthorId,request.Body)
	return c.JSON(comment)
}
