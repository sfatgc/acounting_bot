package accounting_bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processMessageText(ctx context.Context, message *tgbotapi.Message) string {
	return message.Text
}
