package repository

import (
	"context"
	"playoo/internal/domain/entity"
	"gorm.io/gorm"
	"errors"
	"time"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type EventRepository interface{
	Create(c context.Context, event entity.Event) error
	Save(c context.Context ,event entity.Event) error
	Delete(c context.Context, event entity.Event) error
	FindById(c context.Context, id string) (*entity.Event, error)
	FindByAmount(c context.Context, amount int) ([]entity.Event,error)
	FindUpcoming(c context.Context, time time.Time) ([]entity.Event,error)
}

type eventRepository struct {
	DB 		*gorm.DB
	Redis   *redis.Client
}

func NewEventRepository(db *gorm.DB, redis *redis.Client) EventRepository{
	return &eventRepository{
		DB: db,
		Redis: redis,
	}
}

func (er *eventRepository) Create(c context.Context ,event entity.Event) error{
	if err:=er.DB.WithContext(c).Create(&event).Error;err!=nil{
		return err
	}
	eventdata,err:=json.Marshal(event)
	if err!=nil{
		return err
	}
	if er.Redis!=nil{
		if err:=er.Redis.Set(c,event.Id.String(),eventdata,time.Until(event.Time)).Err();err!=nil{
			return err
		}
	}
	return nil
} 

func (er *eventRepository) Save(c context.Context ,event entity.Event) error{
	if err:=er.DB.WithContext(c).Save(&event).Error;err!=nil{
		return err
	}
	if er.Redis!=nil{
		eventdata,err:=json.Marshal(event)
	if err!=nil{
		return err
	}
	ttl, err := er.Redis.TTL(c, event.Id.String()).Result()
    if err != nil {
        return err
    }
	if err:=er.Redis.Set(c,event.Id.String(),eventdata,ttl).Err();err!=nil{
		return err
	}
	}
	return nil
}

func (er *eventRepository) Delete(c context.Context, event entity.Event) error{
	if err:=er.DB.WithContext(c).Delete(&event).Error;err!=nil{
		return err
	}
	return nil
}

func (er *eventRepository) FindUpcoming(c context.Context, time time.Time) ([]entity.Event,error){
	events:=[]entity.Event{}
	if err:=er.DB.WithContext(c).Find(&events,"time <= ?",time).Error;err!=nil{
		return nil,err
	}
	return events,nil
}

func (er *eventRepository) FindById(c context.Context, id string) (*entity.Event, error){
	event:=entity.Event{}
	if er.Redis!=nil{
		eventdata,err:=er.Redis.Get(c,id).Result()
		if err!=nil{
			if err == redis.Nil{
				if err:=er.DB.WithContext(c).First(&event,"id = ?",id).Error;err!=nil{
					return nil,err
				}
				eventdata,err:=json.Marshal(event)
				if err!=nil{
					return nil,err
				}
				if err:=er.Redis.Set(c,id,eventdata,time.Until(event.Time)).Err();err!=nil{
					return nil,err
				}
			}else{
				if err:=er.DB.WithContext(c).First(&event,"id = ?",id).Error;err!=nil{
					return nil,err
				}
			}
		}else{
			if err:=json.Unmarshal([]byte(eventdata),&event);err!=nil{
				return nil,err
			}
		}
	}else{
		if err:=er.DB.WithContext(c).First(&event,"id = ?",id).Error;err!=nil{
			return nil,err
		}
	}	
	return &event,nil
}

func (er *eventRepository) FindByAmount(c context.Context, amount int) ([]entity.Event,error){
	events:=[]entity.Event{}
	if amount < -1{
		return nil, errors.New("incorrect amount")
	}
	if err:=er.DB.WithContext(c).Limit(amount).Find(&events).Error;err!=nil{
		return nil,err
	}
	return events,nil
}

