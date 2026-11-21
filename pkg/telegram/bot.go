package telegram

import (
	"context"
	"os"
	"os/signal"

	"github.com/DuraCHYo/catfacts-tg-bot/pkg/api"
	"github.com/DuraCHYo/catfacts-tg-bot/pkg/translation"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func SendBotMessage(token string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		// bot.WithSkipGetMe(),
	}

	b, err := bot.New(token, opts...)
	if nil != err {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/fact", bot.MatchTypeExact, factHandler)

	b.Start(ctx)
}

func factHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := "Не удалось получить факт, попробуй позже"

	fact, err := api.GetRandomFact()
	if err == nil {
		translated, err := translation.TranslateText(fact)
		if err == nil {
			text = translated
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Напиши /fact",
	})
}
