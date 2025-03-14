package shedulers

import (
	"context"
	"log"
	"playoo/internal/bot"
	"playoo/internal/domain/service"
	"time"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type SheduleEvents struct {
	NoticeService service.NoticeService
	EventService  service.EventService
	UserService   service.UserService
	Logger        *logrus.Logger
	Bot           *bot.Bot
}

func (sheduleEvents *SheduleEvents) SetupSheduleEvents(stop chan struct{}) {
	cr := cron.New()
	cr.AddFunc("@every 1m", func() {
		now := time.Now()
		upcoming, err := sheduleEvents.EventService.FindUpcoming(context.Background(), now.Add(10*time.Minute))
		if err != nil {
			sheduleEvents.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range upcoming {
			if !event.NotifiedPre {
				premsg := "cобытие " + event.Body + " начнется через 10 минут!"
				sheduleEvents.NoticeService.CreateNotice(context.Background(), event, premsg)
				if err := sheduleEvents.Bot.SendMsg(event, premsg); err != nil {
					sheduleEvents.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
				}
				sheduleEvents.Logger.Infof("уведомление о предстоящем событии %v отправлено в %v", event.Body, time.Now())
				event.NotifiedPre = true
				if err := sheduleEvents.EventService.Save(context.Background(), event); err != nil {
					sheduleEvents.Logger.WithError(err).Errorf("failed to save event: %v", err)
				}
			}
		}
		current, err := sheduleEvents.EventService.FindUpcoming(context.Background(), now.Add(1*time.Minute))
		if err != nil {
			sheduleEvents.Logger.WithError(err).Errorf("failed to fetch upcoming events: %v", err)
		}
		for _, event := range current {
			curmsg := "cобытие " + event.Body + " началось!"
			sheduleEvents.NoticeService.CreateNotice(context.Background(), event, curmsg)
			if err := sheduleEvents.Bot.SendMsg(event, curmsg); err != nil {
				sheduleEvents.Logger.WithError(err).Errorf("error to send message to bot: %v", err)
			}
			log.Printf("уведомление о начале события %v отправлено в %v", event.Body, time.Now())
			if err := sheduleEvents.EventService.Delete(context.Background(), event); err != nil {
				sheduleEvents.Logger.WithError(err).Errorf("failed to delete event: %v", err)
			}
			for _, id := range event.Members {
				if err := sheduleEvents.UserService.UpdateEvents(context.Background(), id, event.Id.String()); err != nil {
					sheduleEvents.Logger.WithError(err).Errorf("failed to update event: %v", err)
				}
			}
		}

	})
	cr.Start()
	<-stop
	log.Println("stopping scheduleNotify")
	cr.Stop()
}
