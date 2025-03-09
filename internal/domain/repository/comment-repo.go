package repository

import (
	"context"
	"playoo/internal/domain/entity"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(c context.Context, comment entity.Comment) error
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

