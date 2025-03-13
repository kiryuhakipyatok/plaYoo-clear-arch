package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"playoo/internal/domain/entity"
	"time"
)

type UserRepository interface {
	Create(c context.Context, user entity.User) error
	Save(c context.Context, user entity.User) error
	FindById(c context.Context, id string) (*entity.User, error)
	FindByLogin(c context.Context, login string) (*entity.User, error)
	ExistByLoginOrTg(c context.Context, login, tg string) bool
	FindByAmount(c context.Context, amount int) ([]entity.User, error)
	FindByTg(c context.Context, tg string) (*entity.User, error)
}

type userRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{
		DB:    db,
		Redis: redis,
	}
}

func (ur *userRepository) Create(c context.Context, user entity.User) error {
	if err := ur.DB.WithContext(c).Create(&user).Error; err != nil {
		return err
	}
	if ur.Redis != nil {
		userdata, err := json.Marshal(user)
		if err != nil {
			return err
		}
		ur.Redis.Set(c, user.Id.String(), userdata, time.Hour*24)
	}
	return nil
}

func (ur *userRepository) Save(c context.Context, user entity.User) error {
	if err := ur.DB.WithContext(c).Save(&user).Error; err != nil {
		return err
	}
	if ur.Redis != nil {
		userdata, err := json.Marshal(user)
		if err != nil {
			return err
		}
		if err := ur.Redis.Set(c, user.Id.String(), userdata, time.Hour*24).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (ur *userRepository) FindById(c context.Context, id string) (*entity.User, error) {
	user := entity.User{}
	if ur.Redis != nil {
		userdata, err := ur.Redis.Get(c, id).Result()
		if err != nil {
			if err == redis.Nil {
				if err := ur.DB.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
					return nil, err
				}
				userdata, err := json.Marshal(user)
				if err != nil {
					return nil, err
				}
				if err := ur.Redis.Set(c, id, userdata, time.Hour*24).Err(); err != nil {
					return nil, err
				}
			} else {
				if err := ur.DB.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
					return nil, err
				}
			}
		} else {
			if err := json.Unmarshal([]byte(userdata), &user); err != nil {
				return nil, err
			}
		}
	} else {
		if err := ur.DB.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (ur *userRepository) FindByLogin(c context.Context, login string) (*entity.User, error) {
	user := entity.User{}
	if err := ur.DB.WithContext(c).First(&user, "login = ?", login).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindByTg(c context.Context, tg string) (*entity.User, error) {
	user := entity.User{}
	if err := ur.DB.WithContext(c).First(&user, "tg = ?", tg).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindByAmount(c context.Context, amount int) ([]entity.User, error) {
	users := []entity.User{}
	if amount < -1 {
		return nil, errors.New("incorrect amount")
	}
	if err := ur.DB.WithContext(c).Limit(amount).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *userRepository) ExistByLoginOrTg(c context.Context, login, tg string) bool {
	user := entity.User{}
	if err := ur.DB.WithContext(c).First(&user, "login = ? or tg = ?", login, tg).Error; err != nil {
		return false
	}
	return true
}
