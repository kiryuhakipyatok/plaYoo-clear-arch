package service

import (
	"context"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"time"
	"github.com/google/uuid"
)

type CommentService interface {
	AddCommentToUser(c context.Context,id,rid,body string) (*entity.Comment,error)
	AddCommentToEvent(c context.Context,id,rid,body string) (*entity.Comment,error)
	AddCommentToNews(c context.Context,id,rid,body string) (*entity.Comment,error)
}

type commentService struct {
	CommentRepository repository.CommentRepository
	UserRepository repository.UserRepository
	EventRepository repository.EventRepository
	NewsRepository repository.NewsRepository
}

func NewCommentService(
	commentRepository repository.CommentRepository,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	newsRepository repository.NewsRepository) CommentService{
		return &commentService{
			CommentRepository: commentRepository,
			UserRepository: userRepository,
			EventRepository: eventRepository,
			NewsRepository: newsRepository,
		}
}

func (nc commentService) AddCommentToUser(c context.Context,id,rid,body string) (*entity.Comment,error){
	user,err:=nc.UserRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	reciever,err:=nc.UserRepository.FindById(c,rid)
	if err!=nil{
		return nil,err
	}
	comment:=entity.Comment{
		Id:uuid.New(),
		AuthorId: user.Id,
		AuthorName: user.Login,
		AuthorAvatar: user.Avatar,
		Body: body,
		Receiver: reciever.Id,
		Time: time.Now(),
	}
	if err:=nc.CommentRepository.Create(c,comment);err!=nil{
		return nil,err
	}
	reciever.Comments=append(reciever.Comments, comment.Id.String())
	if err:=nc.UserRepository.Save(c,*reciever);err!=nil{
		return nil,err
	}
	return &comment,nil
}

func (nc commentService) AddCommentToEvent(c context.Context,id,rid,body string) (*entity.Comment,error){
	user,err:=nc.UserRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	event,err:=nc.EventRepository.FindById(c,rid)
	if err!=nil{
		return nil,err
	}
	comment:=entity.Comment{
		Id:uuid.New(),
		AuthorId: user.Id,
		AuthorName: user.Login,
		AuthorAvatar: user.Avatar,
		Body: body,
		Receiver: event.Id,
		Time: time.Now(),
	}
	if err:=nc.CommentRepository.Create(c,comment);err!=nil{
		return nil,err
	}
	event.Comments=append(event.Comments, comment.Id.String())
	if err:=nc.EventRepository.Save(c,*event);err!=nil{
		return nil,err
	}
	return &comment,nil
}

func (nc commentService) AddCommentToNews(c context.Context,id,rid,body string) (*entity.Comment,error){
	user,err:=nc.UserRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	news,err:=nc.NewsRepository.FindById(c,rid)
	if err!=nil{
		return nil,err
	}
	comment:=entity.Comment{
		Id:uuid.New(),
		AuthorId: user.Id,
		AuthorName: user.Login,
		AuthorAvatar: user.Avatar,
		Body: body,
		Receiver: news.Id,
		Time: time.Now(),
	}
	if err:=nc.CommentRepository.Create(c,comment);err!=nil{
		return nil,err
	}
	news.Comments=append(news.Comments, comment.Id.String())
	if err:=nc.NewsRepository.Save(c,*news);err!=nil{
		return nil,err
	}
	return &comment,nil
}