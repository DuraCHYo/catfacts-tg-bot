package main

import (
	"log"

	"github.com/DuraCHYo/catfacts-tg-bot/internal/config"
	"github.com/DuraCHYo/catfacts-tg-bot/pkg/telegram"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("ошибка загрузки .env:", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	telegram.SendBotMessage(cfg.Telegram.BotToken)
}
