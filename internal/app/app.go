package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test/internal/bot"
	"test/internal/config"
	"time"
	"github.com/joho/godotenv"
)



func StartApp() {
	if err := godotenv.Load("../../.env");err != nil {
        log.Fatalf("error loading .env file: %v", err.Error())
    }
	postgres, err := config.ConnectToPostgres()
	if err != nil {
		log.Fatalf("error to connect to postgres: %v", err)
	}
	redis,err:=config.ConnectToRedis()
	if err!=nil{
		log.Printf("error to connect to redis: %v", err)
	}
	closePostgres,err:=postgres.DB()
	if err!=nil{
		log.Fatalf("failed to get postgres to close: %v", err)
	}

	defer func(){
		if err:=closePostgres.Close();err!=nil{
			log.Printf("error to close Postgres: %v",err)
		}else{
			log.Printf("close postgres success")
		}
		if redis!=nil{
			if err:=redis.Close();err!=nil{
				log.Printf("error to close redis: %v",err)
			}else{
				log.Printf("close redis success")
			}
		}
	}()
	app:=config.CreateServer()
	var (
		port	= os.Getenv("PORT")
	)
	cfg:=&config.BootstrapConfig{
		App: app,
		Postgres: postgres,
		Redis: redis,
	}
	quit:=make(chan os.Signal,1)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	stop:=make(chan struct{})
	config.Bootstrap(cfg,stop)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if err := app.Listen("0.0.0.0"+port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	botChan := make(chan *bot.Bot)
	go func(){
		defer wg.Done()
		bot:=config.StartBot(cfg,stop)
		if bot == nil {
			log.Println("failed to create bot")
			return
		}
		log.Println("bot created successfully")
		botChan<-bot
		close(botChan)
	}()
	go func(){
		defer wg.Done()
		bot:=<-botChan
		if bot == nil {
			log.Println("bot is nil, cannot start scheduler")
			return
		}
		log.Println("starting scheduler with bot")
		config.StartShedule(cfg,stop,bot)
	}()
	<-quit
	log.Println("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	if err:=app.ShutdownWithContext(ctx);err!=nil{
		log.Fatalf("server forced to shutdown: %v", err)
	}
	close(stop)
	wg.Wait()
	log.Println("server stopped")
}
