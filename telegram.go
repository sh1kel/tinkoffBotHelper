package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func telegramInit(telegramConfig TelegramConfig) {
	bot, err := tgbotapi.NewBotAPI(telegramConfig.Token)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		for _, user := range telegramConfig.AllowedUsers {
			if user == update.Message.From.UserName {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
