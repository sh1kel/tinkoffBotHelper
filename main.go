package main

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func errorHandle(err error) error {
	if err == nil {
		return nil
	}

	if tradingErr, ok := err.(sdk.TradingError); ok {
		if tradingErr.InvalidTokenSpace() {
			tradingErr.Hint = "Do you use sandbox token in production environment or vise verse?"
			return tradingErr
		}
	}

	return err
}

func main() {
	config := GetConfig()
	commandChan := make(chan string)
	responseChan := make(chan string)
	var wg sync.WaitGroup

	wg.Add(3)
	log.Println("Starting trade loop...")
	go TradeLoop(config.App, commandChan, responseChan, &wg)
	log.Println("Starting telegram loop...")

	go telegramLoop(config.Telegram, commandChan, responseChan, &wg)

	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	//rest(&config.App.Token)
	log.Println("Waiting...")

	wg.Wait()
}
