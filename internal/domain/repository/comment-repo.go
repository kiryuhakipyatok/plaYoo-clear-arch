package repository

import (
	"context"
	"test/internal/domain/entity"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(c context.Context, comment entity.Comment) error
	// FindForUserByAmount(c context.Context, id string,amount int) ([]entity.Comment,error)
	// FindForEventByAmount(c context.Context, id string,amount int) ([]entity.Comment,error)
	// FindForNewsByAmount(c context.Context, id string,amount int) ([]entity.Comment,error)
}

type commentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository{
	return &commentRepository{
		DB: db,
	}
}

func (cr *commentRepository) Create(c context.Context, comment entity.Comment) error{
	if err:=cr.DB.WithContext(c).Create(&comment).Error;err!=nil{
		return err
	}
	return nil
}

// func (cr *commentRepository) FindForUserByAmount(c context.Context, id string,amount int) ([]entity.Comment,error){
// 	comments:=[]entity.Comment{}
// 	if err:=cr.DB.WithContext(c).Limit(amount).
// }

// func (cr *commentRepository) FindForEventByAmount(c context.Context, id string,amount int) ([]entity.Comment,error){

// }

// func (cr *commentRepository) FindForNewsByAmount(c context.Context, id string,amount int) ([]entity.Comment,error){

// }
