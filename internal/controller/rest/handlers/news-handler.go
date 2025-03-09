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

type NewsHandler struct{
	NewsService service.NewsService
	Validator 	*validator.Validate
	Logger 		*logrus.Logger
}

func NewNewsHandler(newsService service.NewsService,validator *validator.Validate,logger *logrus.Logger) NewsHandler{
	return NewsHandler{
		NewsService: newsService,
		Logger: logger,
		Validator: validator,
	}
}

func (nh NewsHandler) CreateNews(c *fiber.Ctx) error{
	ctx:=c.Context()

	request:=dto.CreateNewsRequest{}
	var err error
	request.Picture,err = c.FormFile("picture")
	if err!=nil{
		nh.Logger.Error("no file uploaded")
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "no file uploaded",
        })
	}
	if err:=c.BodyParser(&request);err!=nil{
		return e.ErrorParse(c,nh.Logger,"request",err)
	}

	news,err:=nh.NewsService.CreateNews(ctx,request.Title,request.Body,request.Link,request.Picture)
	if err!=nil{
		return e.FailedToCreate(c,nh.Logger,"news",err)
	}

	nh.Logger.Infof("news created: %v",news)

	return c.JSON(news)
}

func (nh NewsHandler) GetNewsById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")
	news,err:=nh.NewsService.GetById(ctx,id)
	if err!=nil{
		return e.NotFound(c,nh.Logger,"news",err)
	}
	nh.Logger.Infof("news %v recieved",news)
	return c.JSON(news)
}

func (nh NewsHandler) GetNewsByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Params("amount")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return e.ErrorParse(c,nh.Logger,"amount",err)
	}
	somenews,err:=nh.NewsService.GetByAmount(ctx,amount)
	if err!=nil{
		return e.ErrorFetching(c,nh.Logger,"news",err)
	}
	nh.Logger.Infof("some news %v recieved",somenews)
	return c.JSON(somenews)
}