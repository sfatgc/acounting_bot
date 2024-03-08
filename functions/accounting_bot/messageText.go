package accounting_bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processMessageText(ctx context.Context, message *tgbotapi.Message) string {
	u := ctx.Value(U).(TelegramUser)
	message_count, err := u.getMessageCount(ctx, FIRESTORE_CLIENT)
	if err != nil {
		return fmt.Sprintf("NaN (%s)", err)
	}
	return fmt.Sprintf("You wrote:\n\t```\n%s\n```\nYou sent %d messages", message.Text, message_count)
}
