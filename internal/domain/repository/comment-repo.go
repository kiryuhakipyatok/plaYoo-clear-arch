package repository

import (
	"context"
	"playoo/internal/domain/entity"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(c context.Context, comment entity.Comment) error
	FindById(c context.Context, id string, amount int) ([]entity.Comment,error)
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

func (cr *commentRepository) FindById(c context.Context, id string, amount int) ([]entity.Comment,error){
	comments:=[]entity.Comment{}
	if err:=cr.DB.WithContext(c).Limit(amount).Where("receiver=?",id).Find(&comments).Error;err!=nil{
		return nil,err
	}
	return comments,nil
}