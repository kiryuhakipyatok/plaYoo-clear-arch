package repository

import (
	"context"
	"gorm.io/gorm"
	"playoo/internal/domain/entity"
)

type NewsRepository interface {
	Create(c context.Context, news entity.News) error
	Save(c context.Context, news entity.News) error
	FindById(c context.Context, id string) (*entity.News, error)
	FindByAmount(c context.Context, amount int) ([]entity.News, error)
}

type newsRepository struct {
	DB *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{
		DB: db,
	}
}

func (nr *newsRepository) Create(c context.Context, news entity.News) error {
	if err := nr.DB.WithContext(c).Create(&news).Error; err != nil {
		return err
	}
	return nil
}

func (nr *newsRepository) Save(c context.Context, news entity.News) error {
	if err := nr.DB.WithContext(c).Save(&news).Error; err != nil {
		return err
	}
	return nil
}

func (nr *newsRepository) FindById(c context.Context, id string) (*entity.News, error) {
	news := entity.News{}
	if err := nr.DB.WithContext(c).First(&news, "id=?", id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (nr *newsRepository) FindByAmount(c context.Context, amount int) ([]entity.News, error) {
	somenews := []entity.News{}
	if err := nr.DB.WithContext(c).Limit(amount).Find(&somenews).Error; err != nil {
		return nil, err
	}
	return somenews, nil
}
