package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	App      ApplicationConfig
	Telegram TelegramConfig
}

type ApplicationConfig struct {
	Token string
}

type TelegramConfig struct {
	Token        string
	AllowedUsers []string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetConfig() Config {
	configFile, _ := os.Open("config.json")
	decoder := json.NewDecoder(configFile)
	config := new(Config)
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}
	return *config
}
