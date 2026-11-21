package config

import (
	"fmt"
	"os"
)

type Config struct {
	Telegram TelegramConfig
}

type TelegramConfig struct {
	BotToken string
}

func Load() (Config, error) {
	cfg := Config{
		Telegram: TelegramConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
	}

	if cfg.Telegram.BotToken == "" {
		return Config{}, fmt.Errorf("env TELEGRAM_BOT_TOKEN не задан")
	}

	return cfg, nil
}
