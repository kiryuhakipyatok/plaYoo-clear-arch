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

type UserHandler struct{
	UserService 	 service.UserService
	CommentService   service.CommentService
	Validator 		*validator.Validate
	Logger 			*logrus.Logger
	ErrorHandler	*e.ErrorHandler
}

func NewUserHandler(userService service.UserService,commentService service.CommentService,validator *validator.Validate,logger *logrus.Logger, eh *e.ErrorHandler) UserHandler{
	return UserHandler{
		UserService: userService,
		CommentService: commentService,
		Logger: logger,
		Validator: validator,
		ErrorHandler: eh,
	}
}

func (uh UserHandler) GetUserById(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")
	user,err:=uh.UserService.GetById(ctx,id)
	if err!=nil{
		return uh.ErrorHandler.NotFound(c,"user",err)
	}
	uh.Logger.Infof("user %s recieved",user.Id)
	return c.JSON(user)
}

func (uh UserHandler) GetUsersByAmount(c *fiber.Ctx) error{
	ctx:=c.Context()
	a:=c.Query("amount")	
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return uh.ErrorHandler.ErrorParse(c,"amount",err)
	}
	users,err:=uh.UserService.GetByAmount(ctx,amount)
	if err!=nil{
		return uh.ErrorHandler.ErrorFetching(c,"users",err)
	}
	uh.Logger.Infof("users %v recieved",users)
	return c.JSON(users)
}

func (uh UserHandler) UploadAvatar(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")	
	file,err := c.FormFile("avatar")
	if err!=nil{
		uh.Logger.WithError(err).Error("no file uploaded")
        c.Status(fiber.StatusBadRequest)
        return c.JSON(fiber.Map{
            "error": "no file uploaded",
        })
	}
	if err:=uh.UserService.UploadAvatar(ctx,id,file);err!=nil{
		uh.Logger.WithError(err).Error("failed to upload avatar")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to upload avatar",
        })
	}
	uh.Logger.Infof("avatar %s uploaded",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (uh UserHandler) DeleteAvatar(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Params("id")	
	if err:=uh.UserService.DeleteAvatar(ctx,id);err!=nil{
		uh.Logger.WithError(err).Error("failed to delete avatar")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to delete avatar",
        })
	}
	uh.Logger.Infof("avatar %s deleted",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (uh *UserHandler) AddComment(c *fiber.Ctx) error{
	ctx:=c.Context()
	request:=dto.UserCommentRequest{}
	
	if err:=c.BodyParser(&request);err!=nil{
		return uh.ErrorHandler.ErrorParse(c,"request",err)
	}
	if err:=uh.Validator.Struct(request);err!=nil{

		return uh.ErrorHandler.FailedToValidate(c,err)
	}
	comment,err:=uh.CommentService.AddCommentToUser(ctx,request.AuthorId,request.ReceiverId,request.Body)
	if err!=nil{
		return uh.ErrorHandler.FailedToCreate(c,"comment to user",err)
	}
	uh.Logger.Infof("comment added to user %s by %s: %s",request.ReceiverId,request.AuthorId,request.Body)
	return c.JSON(comment)
}

func (uh *UserHandler) GetComments(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	a:=c.Query("amount")	
	amount,err:=strconv.Atoi(a)
	if err!=nil{
		return uh.ErrorHandler.ErrorParse(c,"amount",err)
	}
	comments,err:=uh.CommentService.GetComments(ctx,id,amount)
	if err!=nil{
		return uh.ErrorHandler.ErrorFetching(c,"user's comments",err)
	}
	uh.Logger.Infof("comments %v received",comments)
	return c.JSON(comments)
}

func(uh *UserHandler) RecordDiscord(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	ds:=c.Query("discord")
	if err:=uh.UserService.RecordDiscord(ctx,id,ds);err!=nil{
		uh.Logger.WithError(err).Error("failed to record discord")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to record discord",
        })
	}
	uh.Logger.Infof("discord %s recorded",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func(uh *UserHandler) EditRating(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	s:=c.Query("stars")
	stars,err:=strconv.Atoi(s)
	if err!=nil{
		return uh.ErrorHandler.ErrorParse(c,"amount",err)
	}
	if err:=uh.UserService.EditRating(ctx,id,stars);err!=nil{
		uh.Logger.WithError(err).Error("failed to edit rating")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to edit rating",
        })
	}
	uh.Logger.Infof("rating %s edit",id)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (uh *UserHandler) Follow(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	login:=c.Query("login")
	if err:=uh.UserService.Follow(ctx,id,login);err!=nil{
		uh.Logger.WithError(err).Error("failed to follow to user")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to follow to user",
        })
	}
	uh.Logger.Infof("user %s follow to %s",id,login)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (uh *UserHandler) Unfollow(c *fiber.Ctx) error{
	ctx:=c.Context()
	id:=c.Query("id")
	login:=c.Query("login")
	if err:=uh.UserService.Unfollow(ctx,id,login);err!=nil{
		uh.Logger.WithError(err).Error("failed to unfollow from user")
        c.Status(fiber.StatusInternalServerError)
        return c.JSON(fiber.Map{
            "error": "failed to unfollow from user",
        })
	}
	uh.Logger.Infof("user %s unfollow from %s",id,login)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}