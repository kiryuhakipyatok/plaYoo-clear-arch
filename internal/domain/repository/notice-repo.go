package repository

import (
	"context"
	"gorm.io/gorm"
	"playoo/internal/domain/entity"
)

type NoticeRepository interface {
	Create(c context.Context, notification entity.Notice) error
	Delete(c context.Context, notification entity.Notice) error
	FindById(c context.Context, id string) (*entity.Notice, error)
}

type noticeRepository struct {
	DB *gorm.DB
}

func NewNoticeRepository(db *gorm.DB) NoticeRepository {
	return &noticeRepository{
		DB: db,
	}
}

func (nr *noticeRepository) Create(c context.Context, notification entity.Notice) error {
	if err := nr.DB.WithContext(c).Create(&notification).Error; err != nil {
		return err
	}
	return nil
}

func (nr *noticeRepository) Delete(c context.Context, notification entity.Notice) error {
	if err := nr.DB.WithContext(c).Delete(&notification).Error; err != nil {
		return err
	}
	return nil
}

func (nr *noticeRepository) FindById(c context.Context, id string) (*entity.Notice, error) {
	notification := entity.Notice{}
	if err := nr.DB.WithContext(c).First(&notification, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}
