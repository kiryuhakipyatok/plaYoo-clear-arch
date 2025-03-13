package bot

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"playoo/internal/domain/entity"
	"playoo/internal/domain/repository"
	"strconv"
	"sync"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	UserRepository repository.UserRepository
}

func CreateBot(stop chan struct{}, userRepository repository.UserRepository) *Bot {
	var tbt = os.Getenv("TG_BOT_TOKEN")
	var err error
	bot, err := tgbotapi.NewBotAPI(tbt)
	if err != nil {
		log.Printf("error creating bot: %v", err)
		return nil
	}
	Bot := Bot{bot: bot, UserRepository: userRepository}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Bot.listenForUpdates(stop)
	}()
	return &Bot
}

func (b *Bot) SendMsg(event entity.Event, msg string) error {
	for _, id := range event.Members {
		user, err := b.UserRepository.FindById(context.Background(), id)
		if err != nil {
			return err
		}
		if user.ChatId != "" {
			chatID, _ := strconv.ParseInt(user.ChatId, 10, 64)
			message := tgbotapi.NewMessage(chatID, msg)
			if _, err := b.bot.Send(message); err != nil {
				log.Printf("failed to send message to user %s: %v", user.Telegram, err)
			}
		}
	}
	return nil
}

func (b *Bot) listenForUpdates(stop chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := b.bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case <-stop:
			log.Println("stopping bot")
			return
		case update := <-updates:
			if update.Message != nil {
				b.handleMessage(update)
			}
		}
	}

}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	user, err := b.UserRepository.FindByTg(context.Background(), username)
	if err != nil {
		log.Printf("user not found: %v", err)
	}

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Да, хочу"),
			tgbotapi.NewKeyboardButton("❌ Нет, не хочу"),
		),
	)

	switch text {
	case "✅ Да, хочу":
		if user.ChatId == "" {
			if err := b.storeChatID(user, chatID); err != nil {
				log.Printf("failed to store chatId: %v", err)
			}
			msg := tgbotapi.NewMessage(chatID, "Теперь вы будете получать уведомления от plaYoo о начале ивентов!")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы уже подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	case "❌ Нет, не хочу":
		if user.ChatId != "" {
			if err := b.removeChatID(user); err != nil {
				log.Printf("failed to remove chatId: %v", err)
			}
			msg := tgbotapi.NewMessage(chatID, "Вы отписались от уведомлений.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		} else {
			msg := tgbotapi.NewMessage(chatID, "Вы не подписаны на уведомления.")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := b.bot.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "Хотите ли вы получать уведомления о начале ивентов, к которым вы присоединились?")
		msg.ReplyMarkup = keyboard
		if _, err := b.bot.Send(msg); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func (b *Bot) storeChatID(user *entity.User, chatID int64) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = strconv.Itoa(int(chatID))
	if err := b.UserRepository.Save(context.Background(), *user); err != nil {
		return err
	}
	return nil
}

func (b *Bot) removeChatID(user *entity.User) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	user.ChatId = ""
	if err := b.UserRepository.Save(context.Background(), *user); err != nil {
		return err
	}
	return nil
}
