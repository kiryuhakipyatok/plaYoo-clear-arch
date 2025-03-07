package service

import (
	"context"
	"test/internal/domain/entity"
	"test/internal/domain/repository"
	"github.com/google/uuid"
	//"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"test/internal/bot"
)

type NoticeService interface {
	CreateNotice(c context.Context,event entity.Event,msg string) error
	DeleteNotice(c context.Context,id,nid string) error
	GetNoticeByAmount(c context.Context,id string, amount int) ([]entity.Notice,error)
	DeleteAllNotifications(c context.Context,id string) error
}

type noticeService struct {
	NoticeRepository repository.NoticeRepository
	EventRepository repository.EventRepository
	UserRepository repository.UserRepository
}

func NewNoticeService(
	noticeRepository repository.NoticeRepository,
	eventRepository repository.EventRepository,
	userRepository repository.UserRepository,) NoticeService{
	return &noticeService{
		NoticeRepository: noticeRepository,
		EventRepository: eventRepository,
		UserRepository: userRepository,
	}
}

func (ns *noticeService) CreateNotice(c context.Context,event entity.Event,msg string) error{
	for _,id:=range event.Members{
		user,err:=ns.UserRepository.FindById(c,id)
		if err!=nil{
			return err
		}
		notice:=entity.Notice{
			Id:uuid.New(),
			EventId: event.Id,
			Body: msg,
		}
		if err:=ns.NoticeRepository.Create(c,notice);err!=nil{
			return err
		}
		user.Notifications = append(user.Notifications, notice.Id.String())
		if err:=ns.UserRepository.Save(c,*user);err!=nil{
			return err
		}
	}
	return nil
}

func (ns *noticeService) DeleteNotice(c context.Context,id,nid string) error{
	notice,err:=ns.NoticeRepository.FindById(c,nid)
	if err!=nil{
		return err
	}
	if err:=ns.NoticeRepository.Delete(c,*notice);err!=nil{
		return err
	}
	user,err:=ns.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	updateNotifications:=make([]string,0,len(user.Notifications))
	for _,n:=range user.Notifications{
		if n!=nid{
			updateNotifications = append(updateNotifications, n)
		}
	}
	user.Notifications = updateNotifications
	if err:=ns.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

func (ns *noticeService) GetNoticeByAmount(c context.Context, id string,amount int) ([]entity.Notice,error){
	user,err:=ns.UserRepository.FindById(c,id)
	if err!=nil{
		return nil,err
	}
	notifications:=[]entity.Notice{}
	for _,n:= range user.Notifications[:amount]{
		nts,err:=ns.NoticeRepository.FindById(c,n)
		if err!=nil{
			return nil,err
		}
		notifications = append(notifications, *nts)
	}
	
	return notifications,nil
}

func (ns *noticeService) DeleteAllNotifications(c context.Context,id string) error{
	user,err:=ns.UserRepository.FindById(c,id)
	if err!=nil{
		return err
	}
	user.Notifications = nil
	if err:=ns.UserRepository.Save(c,*user);err!=nil{
		return err
	}
	return nil
}

