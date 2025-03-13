package service

import (
	"context"
	"github.com/google/uuid"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"time"
)

type CommentService interface {
	AddCommentToUser(c context.Context, id, rid, body string) (*entity.Comment, error)
	AddCommentToEvent(c context.Context, id, rid, body string) (*entity.Comment, error)
	AddCommentToNews(c context.Context, id, rid, body string) (*entity.Comment, error)
	GetComments(c context.Context, id string, amount int) ([]entity.Comment, error)
}

type commentService struct {
	CommentRepository repository.CommentRepository
	UserRepository    repository.UserRepository
	EventRepository   repository.EventRepository
	NewsRepository    repository.NewsRepository
	Transactor        repository.Transactor
}

func NewCommentService(
	commentRepository repository.CommentRepository,
	userRepository repository.UserRepository,
	eventRepository repository.EventRepository,
	newsRepository repository.NewsRepository,
	transactor repository.Transactor) CommentService {
	return &commentService{
		CommentRepository: commentRepository,
		UserRepository:    userRepository,
		EventRepository:   eventRepository,
		NewsRepository:    newsRepository,
		Transactor:        transactor,
	}
}

func (cs commentService) AddCommentToUser(c context.Context, id, rid, body string) (*entity.Comment, error) {
	res, err := cs.Transactor.WithinTransaction(c, func(c context.Context) (any, error) {
		user, err := cs.UserRepository.FindById(c, id)
		if err != nil {
			return nil, err
		}
		reciever, err := cs.UserRepository.FindById(c, rid)
		if err != nil {
			return nil, err
		}
		comment := entity.Comment{
			Id:           uuid.New(),
			AuthorId:     user.Id,
			AuthorName:   user.Login,
			AuthorAvatar: user.Avatar,
			Body:         body,
			Receiver:     reciever.Id,
			Time:         time.Now(),
		}
		if err := cs.CommentRepository.Create(c, comment); err != nil {
			return nil, err
		}
		reciever.Comments = append(reciever.Comments, comment.Id.String())
		if err := cs.UserRepository.Save(c, *reciever); err != nil {
			return nil, err
		}
		return &comment, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*entity.Comment), nil
}

func (cs commentService) AddCommentToEvent(c context.Context, id, rid, body string) (*entity.Comment, error) {
	res, err := cs.Transactor.WithinTransaction(c, func(c context.Context) (any, error) {
		user, err := cs.UserRepository.FindById(c, id)
		if err != nil {
			return nil, err
		}
		event, err := cs.EventRepository.FindById(c, rid)
		if err != nil {
			return nil, err
		}
		comment := entity.Comment{
			Id:           uuid.New(),
			AuthorId:     user.Id,
			AuthorName:   user.Login,
			AuthorAvatar: user.Avatar,
			Body:         body,
			Receiver:     event.Id,
			Time:         time.Now(),
		}
		if err := cs.CommentRepository.Create(c, comment); err != nil {
			return nil, err
		}
		event.Comments = append(event.Comments, comment.Id.String())
		if err := cs.EventRepository.Save(c, *event); err != nil {
			return nil, err
		}
		return &comment, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*entity.Comment), nil
}

func (cs commentService) AddCommentToNews(c context.Context, id, rid, body string) (*entity.Comment, error) {
	res, err := cs.Transactor.WithinTransaction(c, func(c context.Context) (any, error) {
		user, err := cs.UserRepository.FindById(c, id)
		if err != nil {
			return nil, err
		}
		news, err := cs.NewsRepository.FindById(c, rid)
		if err != nil {
			return nil, err
		}
		comment := entity.Comment{
			Id:           uuid.New(),
			AuthorId:     user.Id,
			AuthorName:   user.Login,
			AuthorAvatar: user.Avatar,
			Body:         body,
			Receiver:     news.Id,
			Time:         time.Now(),
		}
		if err := cs.CommentRepository.Create(c, comment); err != nil {
			return nil, err
		}
		news.Comments = append(news.Comments, comment.Id.String())
		if err := cs.NewsRepository.Save(c, *news); err != nil {
			return nil, err
		}
		return &comment, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*entity.Comment), nil
}

func (cs commentService) GetComments(c context.Context, id string, amount int) ([]entity.Comment, error) {
	comments, err := cs.CommentRepository.FindById(c, id, amount)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
