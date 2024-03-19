package accounting_bot

import (
	"net/http"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botRuntime struct {
	user *TelegramUser
	r    *http.Request
	db   *firestore.Client
	tg   *tgbotapi.BotAPI
}

func newRuntime(r *http.Request, db *firestore.Client, tg *tgbotapi.BotAPI) *botRuntime {
	return &botRuntime{nil, r, db, tg}
}
