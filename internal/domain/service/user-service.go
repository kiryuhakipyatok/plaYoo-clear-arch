package service

import (
	"context"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
)

type UserService interface {
	GetById(c context.Context, id string) (*entity.User, error)
	GetByAmount(c context.Context, amount int) ([]entity.User,error)
	UpdateEvents(c context.Context, id,eventid string) error
}

type userService struct{
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService{
	return &userService{
		UserRepository: userRepository,
	}
}	

func (us userService) GetById(c context.Context, id string) (*entity.User, error){
	return us.UserRepository.FindById(c,id)
}

func (us userService) GetByAmount(c context.Context, amount int) ([]entity.User,error){
	return us.UserRepository.FindByAmount(c,amount)
}


func (ur userService) UpdateEvents(c context.Context, id,eventid string) error{
	user,err:=ur.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	updateEvents:=make([]string,0,len(user.Events))
	for _, e := range user.Events {
		if e != eventid{
			updateEvents = append(updateEvents, e)
		}
	}
	user.Events = updateEvents
	if err:=ur.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}
