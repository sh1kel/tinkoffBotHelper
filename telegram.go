package main

import (
	"log"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func telegramLoop(telegramConfig TelegramConfig, commandChan chan string, responseChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	updateChan := make(chan tgbotapi.Update, 10)

	bot, err := tgbotapi.NewBotAPI(telegramConfig.Token)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go ProcessTelegramMessage(bot, commandChan, responseChan, updateChan, wg)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		for _, user := range telegramConfig.AllowedUsers {
			if user == update.Message.From.UserName {
				updateChan <- update
				break
			}
		}
	}
}

func ProcessTelegramMessage(bot *tgbotapi.BotAPI, commandChan chan string, responseChan chan string, updateChan chan tgbotapi.Update, wg *sync.WaitGroup) {
	defer wg.Done()

	for update := range updateChan {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		commandChan <- update.Message.Text
		responseText := <-responseChan
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = "Markdown"

		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
