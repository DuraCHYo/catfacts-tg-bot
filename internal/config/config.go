package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Telegram TelegramConfig
}

type TelegramConfig struct {
	BotToken string
}

func Load() (Config, error) {
	// Настраиваю глобальное логирование при старте бота
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)
	slog.SetDefault(logger)

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
