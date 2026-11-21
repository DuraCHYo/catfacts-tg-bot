package telegram

import (
	"context"
	"log/slog"
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

	languages := NewLanguageStore()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		// Перед каждым хендлером кладём в контекст язык пользователя.
		bot.WithMiddlewares(withUserLanguage(languages)),
	}

	b, err := bot.New(token, opts...)
	if nil != err {
		slog.Error("Ошибка создания соединения бота", "error", err)
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "fact", bot.MatchTypeCommand, factHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "language", bot.MatchTypeCommand, changeLanguageHandler(languages))

	slog.Info("Бот запущен...")
	b.Start(ctx)
	slog.Warn("Бот остановлен...")
}

func factHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	lang := LanguageFromContext(ctx)

	text := localized(lang,
		"Не удалось получить факт, попробуй позже",
		"Could not get a fact, please try again later",
	)

	fact, err := api.GetRandomFact()
	if err == nil {
		if lang == "ru" {
			if translated, err := translation.TranslateText(fact); err == nil {
				text = translated
			}
		} else {
			text = fact
		}
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

// changeLanguageHandler при каждом вызове переключает язык по кругу: ru → en → ru.
func changeLanguageHandler(languages *LanguageStore) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := userIDFromUpdate(update)
		if userID == 0 {
			return
		}

		// Меняем язык в сторе. Следующий update уже получит его
		// из middleware, поэтому выбор действует для всех хендлеров.
		newLang := languages.Next(userID)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: localized(newLang,
				"Язык изменён на русский",
				"Language changed to English",
			),
		})
	}
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := localized(LanguageFromContext(ctx),
		`Бот для получения фактов о котах.
Поддерживает русский и английский языки.
Чтобы получить факт — вызовите /fact.
Чтобы переключить язык — вызовите /language.`,
		`A bot that sends cat facts.
Supports Russian and English.
To get a fact, send /fact.
To switch the language, send /language.`,
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   localized(LanguageFromContext(ctx), "Напиши /help", "Send /help"),
	})
}
