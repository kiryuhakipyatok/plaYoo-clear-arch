package handlers

import (
	"strconv"
	"test/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct{
	NewsService service.NewsService
}

func NewNewsHandler(newsService service.NewsService) NewsHandler{
	return NewsHandler{
		NewsService: newsService,
	}
}

func (nh NewsHandler) CreateNews(c *fiber.Ctx) error{
	ctx:=c.Context()

	title := c.FormValue("title")
    body := c.FormValue("body")
    link := c.FormValue("link")
	picture, err := c.FormFile("picture")

    if err != nil {
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "no file uploaded",
        })
    }

	news,err:=nh.NewsService.CreateNews(ctx,title,body,link,picture)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"error": "error to create news",
		})
	}

	return c.JSON(news)
}

func (nh NewsHandler) GetNewsById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	news,err:=nh.NewsService.GetById(ctx,id)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
        return c.JSON(fiber.Map{
            "error": "news not found",
        })
	}
	return c.JSON(news)
}

func (nh NewsHandler) GetNewsByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "error parse amount",
        })
	}
	somenews,err:=nh.NewsService.GetByAmount(ctx,amount)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
        return c.JSON(fiber.Map{
            "error": "some news not found",
        })
	}
	return c.JSON(somenews)
}