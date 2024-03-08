package accounting_bot

import (
	"context"

	"cloud.google.com/go/firestore"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botRuntime struct {
	user *TelegramUser
	rCtx context.Context
	db   *firestore.Client
	tg   *tgbotapi.BotAPI
}

func newRuntime(ctx context.Context, db *firestore.Client, tg *tgbotapi.BotAPI) *botRuntime {
	return &botRuntime{nil, ctx, db, tg}
}
