package telegram

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const DefaultLanguage = "ru"

// supportedLanguages задаёт порядок циклического переключения.
var supportedLanguages = []string{"ru", "en"}

// languageKey — типизированный ключ контекста, чтобы не пересечься
// со значениями из других пакетов.
type languageKey struct{}

// LanguageStore хранит язык каждого пользователя.
// map защищён мьютексом, потому что update'ы обрабатываются конкурентно.
type LanguageStore struct {
	mu     sync.RWMutex
	byUser map[int64]string
}

func NewLanguageStore() *LanguageStore {
	return &LanguageStore{byUser: make(map[int64]string)}
}

func (s *LanguageStore) Get(userID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if lang, ok := s.byUser[userID]; ok {
		return lang
	}

	return DefaultLanguage
}

// Next переключает язык пользователя на следующий по кругу
// и возвращает новый язык.
func (s *LanguageStore) Next(userID int64) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.byUser[userID]
	if !ok {
		current = DefaultLanguage
	}

	next := supportedLanguages[0]
	for i, lang := range supportedLanguages {
		if lang == current {
			next = supportedLanguages[(i+1)%len(supportedLanguages)]
			break
		}
	}

	s.byUser[userID] = next
	return next
}

// LanguageFromContext возвращает язык, который middleware положил в контекст.
func LanguageFromContext(ctx context.Context) string {
	lang, ok := ctx.Value(languageKey{}).(string)
	if !ok {
		return DefaultLanguage
	}

	return lang
}

// withUserLanguage подставляет язык пользователя в контекст каждого update.
// Вызывается до хендлера, поэтому в любом хендлере доступен актуальный язык.
func withUserLanguage(store *LanguageStore) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if userID := userIDFromUpdate(update); userID != 0 {
				ctx = context.WithValue(ctx, languageKey{}, store.Get(userID))
			}

			next(ctx, b, update)
		}
	}
}

func userIDFromUpdate(update *models.Update) int64 {
	switch {
	case update.Message != nil && update.Message.From != nil:
		return update.Message.From.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.ID
	case update.InlineQuery != nil && update.InlineQuery.From != nil:
		return update.InlineQuery.From.ID
	default:
		return 0
	}
}

// localized выбирает строку по языку пользователя.
func localized(lang, ru, en string) string {
	if lang == "ru" {
		return ru
	}

	return en
}
