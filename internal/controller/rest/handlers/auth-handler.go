package handlers

import (
	"test/internal/domain/service"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"os"
	// "github.com/joho/godotenv"
	"log"
)

type AuthHandler struct{
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler{
	return AuthHandler{
		AuthService: authService,
	}
}

func (ah AuthHandler) Register(c *fiber.Ctx) error{
	ctx:=c.Context()
	var request struct{
		Login 		string 	`json:"login"`
		Telegram 	string 	`json:"telegram"`
		Password 	string 	`json:"password"`
	}
	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"bad request" + err.Error(),
		})
	}
	user,err:=ah.AuthService.Register(ctx,request.Login,request.Telegram,request.Password)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"failed register" + err.Error(),
		})
	}
	return c.JSON(user)
}

func (ah AuthHandler) Login(c *fiber.Ctx) error{
	ctx:=c.Context()

	var request struct{
		Login 		string 	`json:"login"`
		Password 	string 	`json:"password"`
	}

	if err:=c.BodyParser(&request);err!=nil{
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error":"bad request" + err.Error(),
		})
	}

	token,err:=ah.AuthService.GetTokenForLogin(ctx,request.Login,request.Password)

	if err != nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"failed login" + err.Error(),
		})
	}

	cookie:=fiber.Cookie{
		Name:"jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (ah AuthHandler) Logout(c *fiber.Ctx) error{
	cookie:=fiber.Cookie{
		Name:"jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":"success",
	})
}

func (ah AuthHandler) GetLoggedUser(c *fiber.Ctx) error{
	// if err := godotenv.Load("../../.env");err != nil {
    //     log.Fatalf("error loading .env file when start server: %v", err.Error())
    // }
	secret:=os.Getenv("SECRET")
	if secret == ""{
		log.Fatal("error secret .env value  is empty")
	}
	ctx:=c.Context()
	cookie:=c.Cookies("jwt")
	token,err:=jwt.ParseWithClaims(cookie,&jwt.StandardClaims{},func(t *jwt.Token) (interface{}, error) {
		return []byte(secret),nil
	})
	if err!=nil{
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"error":"unauthenticated",
		})
	}
	claims:=token.Claims.(*jwt.StandardClaims)
	user,err:=ah.AuthService.GetUserByClaims(ctx,claims.Issuer)
	if err!=nil{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"error":"user not found",
		})
	}
	return c.JSON(user)
}

