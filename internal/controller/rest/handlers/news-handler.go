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

type NewsHandler struct {
	NewsService    service.NewsService
	CommentService service.CommentService
	Validator      *validator.Validate
	Logger         *logrus.Logger
	ErrorHandler   *e.ErrorHandler
}

func NewNewsHandler(newsService service.NewsService, commentService service.CommentService, validator *validator.Validate, logger *logrus.Logger, eh *e.ErrorHandler) NewsHandler {
	return NewsHandler{
		NewsService:    newsService,
		CommentService: commentService,
		Logger:         logger,
		Validator:      validator,
		ErrorHandler:   eh,
	}
}

func (nh *NewsHandler) CreateNews(c *fiber.Ctx) error {
	ctx := c.Context()

	request := dto.CreateNewsRequest{}
	var err error
	request.Picture, err = c.FormFile("picture")
	if err != nil {
		nh.Logger.Error("no file uploaded")
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": "no file uploaded",
		})
	}
	if err := c.BodyParser(&request); err != nil {
		return nh.ErrorHandler.ErrorParse(c, "request", err)
	}

	news, err := nh.NewsService.CreateNews(ctx, request.Title, request.Body, request.Link, request.Picture)
	if err != nil {
		return nh.ErrorHandler.FailedToCreate(c, "news", err)
	}

	nh.Logger.Infof("news created: %s", news.Id)
	response := dto.NewsResponse{
		Id:    news.Id,
		Title: news.Title,
	}
	return c.JSON(response)
}

func (nh *NewsHandler) GetNewsById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	news, err := nh.NewsService.GetById(ctx, id)
	if err != nil {
		return nh.ErrorHandler.NotFound(c, "news", err)
	}
	nh.Logger.Infof("news %s recieved", news.Id)
	return c.JSON(news)
}

func (nh *NewsHandler) GetNewsByAmount(c *fiber.Ctx) error {
	ctx := c.Context()
	a := c.Query("amount")
	amount, err := strconv.Atoi(a)
	if err != nil {
		return nh.ErrorHandler.ErrorParse(c, "amount", err)
	}
	somenews, err := nh.NewsService.GetByAmount(ctx, amount)
	if err != nil {
		return nh.ErrorHandler.ErrorFetching(c, "news", err)
	}
	nh.Logger.Infof("some news %v recieved", somenews)
	return c.JSON(somenews)
}

func (nh *NewsHandler) AddComment(c *fiber.Ctx) error {
	ctx := c.Context()
	request := dto.NewsCommentRequest{}

	if err := c.BodyParser(&request); err != nil {
		return nh.ErrorHandler.ErrorParse(c, "request", err)
	}
	if err := nh.Validator.Struct(request); err != nil {
		return nh.ErrorHandler.FailedToValidate(c, err)
	}
	comment, err := nh.CommentService.AddCommentToUser(ctx, request.AuthorId, request.ReceiverId, request.Body)
	if err != nil {
		return nh.ErrorHandler.FailedToCreate(c, "comment to news", err)
	}
	nh.Logger.Infof("comment added to news %s by %s: %s", request.ReceiverId, request.AuthorId, request.Body)
	return c.JSON(comment)
}

func (nh *NewsHandler) GetComments(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Query("id")
	a := c.Query("amount")
	amount, err := strconv.Atoi(a)
	if err != nil {
		return nh.ErrorHandler.ErrorParse(c, "amount", err)
	}
	comments, err := nh.CommentService.GetComments(ctx, id, amount)
	if err != nil {
		return nh.ErrorHandler.ErrorFetching(c, "news's comments", err)
	}
	nh.Logger.Infof("comments %v received", comments)
	return c.JSON(comments)
}
