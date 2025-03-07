package handlers

import (
	"strconv"
	"test/internal/domain/service"
	"github.com/gofiber/fiber/v2"
)

type NoticeHandler struct{
	NoticeService service.NoticeService
}

func NewNoticeHandler(noticeService service.NoticeService) NoticeHandler{
	return NoticeHandler{
		NoticeService: noticeService,
	}
}

func (nh NoticeHandler) DeleteNotice(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	nid:=c.Query("notice")
	if err:=nh.NoticeService.DeleteNotice(ctx,id,nid);err!=nil{
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to delete notice",
        })
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (nh NoticeHandler) GetNotifications(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")
	id:=c.Query("id")
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "error parse amount",
        })
	}
	notifications,err:=nh.NoticeService.GetNoticeByAmount(ctx,id,amount)
	if err!=nil{
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to fetch notifications",
        })
	}
	return c.JSON(notifications)
}

func (nh NoticeHandler) DeleteAllNotifications(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	if err:=nh.NoticeService.DeleteAllNotifications(ctx,id);err!=nil{
		c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to delete all notifications",
        })
	}
	return c.JSON(fiber.Map{
		"message":"success",
	})
}
