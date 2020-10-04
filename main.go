package main

import (
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"math/rand"
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
	telegramInit(config.Telegram)
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	//rest(&config.App.Token)

}
