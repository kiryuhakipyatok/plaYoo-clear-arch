package handlers

import (
	"strconv"
	"playoo/internal/domain/service"
	"github.com/gofiber/fiber/v2"
	e "playoo/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type NoticeHandler struct{
	NoticeService 	 service.NoticeService
	Validator 		*validator.Validate
	Logger 			*logrus.Logger
	ErrorHandler 	*e.ErrorHandler
}

func NewNoticeHandler(noticeService service.NoticeService,validator *validator.Validate,logger *logrus.Logger,eh *e.ErrorHandler) NoticeHandler{
	return NoticeHandler{
		NoticeService: noticeService,
		Logger: logger,
		Validator: validator,
		ErrorHandler: eh,
	}
}

func (nh NoticeHandler) DeleteNotice(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	nid:=c.Query("notice")
	if err:=nh.NoticeService.DeleteNotice(ctx,id,nid);err!=nil{
		return nh.ErrorHandler.FailedToDelete(c,"notification",err)
	}
	nh.Logger.Infof("notification %s deleted",nid)
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
		return nh.ErrorHandler.ErrorParse(c,"amount",err)
	}
	notifications,err:=nh.NoticeService.GetNoticeByAmount(ctx,id,amount)
	if err!=nil{
		return nh.ErrorHandler.ErrorFetching(c,"notifications",err)
	}
	nh.Logger.Infof("notifications %v recieved",notifications)
	return c.JSON(notifications)
}

func (nh NoticeHandler) DeleteAllNotifications(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")
	if err:=nh.NoticeService.DeleteAllNotifications(ctx,id);err!=nil{
		return nh.ErrorHandler.FailedToDelete(c,"all notification",err)
	}
	nh.Logger.Infof("user's %s notifications deleted",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}
