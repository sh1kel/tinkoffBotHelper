package main

import (
	"encoding/json"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"log"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	App ApplicationConfig
}

type ApplicationConfig struct {
	Token string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
	rand.Seed(time.Now().UnixNano()) // инициируем Seed рандома для функции requestID
	configFile, _ := os.Open("config.json")
	decoder := json.NewDecoder(configFile)
	config := new(Config)
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	rest(&config.App.Token)

	//stream()
}
