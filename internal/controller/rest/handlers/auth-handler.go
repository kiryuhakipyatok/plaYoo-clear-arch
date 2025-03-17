package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"os"
	"playoo/internal/domain/service"
	"playoo/internal/dto"
	e "playoo/pkg/errors"
	"time"
)

type AuthHandler struct {
	AuthService  service.AuthService
	Validator    *validator.Validate
	ErrorHandler *e.ErrorHandler
	Logger       *logrus.Logger
}

func NewAuthHandler(authService service.AuthService, validator *validator.Validate, logger *logrus.Logger, eh *e.ErrorHandler) AuthHandler {
	return AuthHandler{
		AuthService:  authService,
		Validator:    validator,
		ErrorHandler: eh,
		Logger:       logger,
	}
}

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	ctx := c.Context()
	request := dto.RegisterRequest{}
	if err := c.BodyParser(&request); err != nil {
		return ah.ErrorHandler.ErrorParse(c, "request", err)
	}
	if err := ah.Validator.Struct(request); err != nil {
		return ah.ErrorHandler.FailedToValidate(c, err)
	}
	user, err := ah.AuthService.Register(ctx, request.Login, request.Telegram, request.Password)
	if err != nil {
		ah.Logger.WithError(err).Error("failed register")
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "failed register: " + err.Error(),
		})
	}
	ah.Logger.Infof("user registered: %s", user.Id)
	respone := dto.RegisterResponse{
		Id:       user.Id,
		Login:    user.Login,
		Telegram: user.Telegram,
	}
	return c.JSON(respone)
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	request := dto.LoginRequest{}

	if err := c.BodyParser(&request); err != nil {
		return ah.ErrorHandler.ErrorParse(c, "request", err)
	}
	if err := ah.Validator.Struct(request); err != nil {
		return ah.ErrorHandler.FailedToValidate(c, err)
	}
	token, err := ah.AuthService.GetTokenForLogin(ctx, request.Login, request.Password)

	if err != nil {
		ah.Logger.WithError(err).Error("failed login")
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error": "failed login: " + err.Error(),
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (ah *AuthHandler) Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func (ah *AuthHandler) GetLoggedUser(c *fiber.Ctx) error {
	secret := os.Getenv("SECRET")
	if secret == "" {
		ah.Logger.Fatal("error secret .env value  is empty")
	}
	ctx := c.Context()
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		ah.Logger.WithError(err).Info("unauthenticated")
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	user, err := ah.AuthService.GetUserByClaims(ctx, claims.Issuer)
	if err != nil {
		return ah.ErrorHandler.NotFound(c, "user", err)
	}

	ah.Logger.Infof("logged in user: %s", user.Id)

	return c.JSON(user)
}
