package config

import (
	"fmt"
	"log"
	"test/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"github.com/joho/godotenv"
	"os"
)

func ConnectToPostgres() (*gorm.DB,error){
	var (
		db	= 	os.Getenv("POSTGRES_DB")
		password	= 	os.Getenv("POSTGRES_PASSWORD")
		user		= 	os.Getenv("POSTGRES_USER")
		port		= 	os.Getenv("PGPORT")
		host		= 	os.Getenv("PGHOST")
	)
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Minsk",host,user,password,db,port)
	connection,err:= gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		return nil,err
	}
	connection.AutoMigrate(
		&entity.User{},
		&entity.Event{},
		&entity.Comment{},
		&entity.Game{},
		&entity.News{},
		&entity.Notice{},
	)
	log.Printf("connect to postgress successfully")
	return connection,nil
}