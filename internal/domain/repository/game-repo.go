package repository

import (
	"context"
	"test/internal/domain/entity"
	"gorm.io/gorm"
)

type GameRepository interface{
	Create(c context.Context,game entity.Game) error
	Save(c context.Context,game entity.Game) error
	FindByName(c context.Context, name string) (*entity.Game,error)
	FindByAmount(c context.Context, amount int) ([]entity.Game,error)
}

type gameRepository struct {
	DB *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository{
	return &gameRepository{
		DB: db,
	}
}

func (gr *gameRepository) Create(c context.Context,game entity.Game) error{
	if err:=gr.DB.WithContext(c).Create(&game).Error;err!=nil{
		return err
	}
	return nil
}
func (gr *gameRepository) Save(c context.Context,game entity.Game) error{
	if err:=gr.DB.WithContext(c).Save(&game).Error;err!=nil{
		return err
	}
	return nil
}
func (gr *gameRepository) FindByName(c context.Context, name string) (*entity.Game,error){
	game:=entity.Game{}
	if err:=gr.DB.WithContext(c).First(&game,"name=?",name).Error;err!=nil{
		return nil,err
	}
	return &game,nil
}
func (gr *gameRepository) FindByAmount(c context.Context, amount int) ([]entity.Game,error){
	games:=[]entity.Game{}
	if err:=gr.DB.WithContext(c).Limit(amount).Find(&games).Error;err!=nil{
		return nil,err
	}
	return games,nil
}


