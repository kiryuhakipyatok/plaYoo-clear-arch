package handlers

import (
	"test/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct{
	CommentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) CommentHandler{
	return CommentHandler{
		CommentService: commentService,
	}
}

func (ch CommentHandler) AddCommentToUser(c *fiber.Ctx) error{
	ctx:=c.Context()
	var request struct{
		AuthorId string `json:"author-id"`
		ReceiverId string `json:"reciever-id"`
		Body string `json:"body"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"error parse request",
		})
	}
	comment,err:=ch.CommentService.AddCommentToUser(ctx,request.AuthorId,request.ReceiverId,request.Body)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error create comment to user",
		})
	}
	return c.JSON(comment)
}

func (ch CommentHandler) AddCommentToEvent(c *fiber.Ctx) error{
	ctx:=c.Context()
	var request struct{
		AuthorId string `json:"author-id"`
		EventId string `json:"event-id"`
		Body string `json:"body"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"error parse request",
		})
	}
	comment,err:=ch.CommentService.AddCommentToEvent(ctx,request.AuthorId,request.EventId,request.Body)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error create comment to event",
		})
	}
	return c.JSON(comment)
}

func (ch CommentHandler) AddCommentToNews(c *fiber.Ctx) error{
	ctx:=c.Context()
	var request struct{
		AuthorId string `json:"author-id"`
		NewsId string `json:"news-id"`
		Body string `json:"body"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"error parse request",
		})
	}
	comment,err:=ch.CommentService.AddCommentToUser(ctx,request.AuthorId,request.NewsId,request.Body)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error":"error create comment to news",
		})
	}
	return c.JSON(comment)
}
