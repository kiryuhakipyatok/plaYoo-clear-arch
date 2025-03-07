package service

import (
	"context"
	"test/internal/domain/entity"
	"test/internal/domain/repository"
	"time"
	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(c context.Context, id,body,name string, max,minute int) (*entity.Event,error)
	GetById(c context.Context, id string) (*entity.Event, error)
	GetByAmount(c context.Context, amount int) ([]entity.Event,error)
	FindUpcoming(c context.Context, time time.Time) ([]entity.Event,error)
	Save(c context.Context ,event entity.Event) error
	Delete(c context.Context, event entity.Event) error
}

type eventService struct {
	EventRepository repository.EventRepository
	UserRepository  repository.UserRepository 
	GameRepository 	repository.GameRepository
}

func NewEventService(
	eventRepository repository.EventRepository,
	userRepository repository.UserRepository,
	gameRepository repository.GameRepository)  EventService{
	return &eventService{
		EventRepository: eventRepository,
		UserRepository: userRepository,
		GameRepository: gameRepository,
	}
}

func (es eventService) CreateEvent(c context.Context, id,body,name string, max,minute int) (*entity.Event,error){
	user,err:=es.UserRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	game,err:=es.GameRepository.FindByName(c,name)
	if err!=nil{
		return nil,err
	}
	event:=entity.Event{
		Id: uuid.New(),
		AuthorId: user.Id,
		Body: body,
		Game: game.Name,
		Max: max,
		Time: time.Now().Add(time.Duration(minute)*time.Minute),
	}
	if minute < 10{
		event.NotifiedPre = true
	}
	event.Members = append(event.Members, user.Id.String())
	user.Events = append(user.Events, event.Id.String())
	game.NumberOfEvents ++
	if err:=es.EventRepository.Create(c,event);err!=nil{
		return nil,err
	}
	if err:=es.GameRepository.Save(c,*game);err!=nil{
		return nil,err
	}
	if err:=es.UserRepository.Save(c,*user);err!=nil{
		return nil,err
	}
	return &event,nil
}

func (es eventService) GetById(c context.Context, id string) (*entity.Event, error){
	event,err:=es.EventRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	return event,nil
}

func (es eventService) GetByAmount(c context.Context, amount int) ([]entity.Event,error){
	events,err:=es.EventRepository.FindByAmount(c,amount)
	if err!=nil{
		return nil,err
	}
	return events,nil
}

func (es eventService) FindUpcoming(c context.Context, time time.Time) ([]entity.Event,error){
	events,err:=es.EventRepository.FindUpcoming(c,time)
	if err!=nil{
		return nil,err
	}
	return events,nil
}

func (es eventService) Save(c context.Context ,event entity.Event) error{
	if err:=es.EventRepository.Save(c,event);err!=nil{
		return err
	}
	return nil
}

func (es eventService) Delete(c context.Context, event entity.Event) error{
	if err:=es.EventRepository.Delete(c,event);err!=nil{
		return err
	}
	return nil
}