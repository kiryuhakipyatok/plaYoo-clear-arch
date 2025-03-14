package service

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"time"
)

type AuthService interface {
	Register(c context.Context, login, tg, password string) (*entity.User, error)
	GetTokenForLogin(c context.Context, login, password string) (string, error)
	GetUserByClaims(c context.Context, claims string) (*entity.User, error)
}

type authService struct {
	UserRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		UserRepository: userRepository,
	}
}

func (as authService) Register(c context.Context, login, tg, password string) (*entity.User, error) {
	if as.UserRepository.ExistByLoginOrTg(c, login, tg) {
		return nil, errors.New("user with same login or tg alredy exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		Id:       uuid.New(),
		Login:    login,
		Telegram: tg,
		Password: string(hashPassword),
	}
	if err := as.UserRepository.Create(c, user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (as authService) GetTokenForLogin(c context.Context, login, password string) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("error secret .env value  is empty")
	}
	user, err := as.UserRepository.FindByLogin(c, login)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as authService) GetUserByClaims(c context.Context, claims string) (*entity.User, error) {
	user, err := as.UserRepository.FindById(c, claims)
	if err != nil {
		return nil, err
	}
	return user, nil
}
