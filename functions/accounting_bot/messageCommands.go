package accounting_bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processMessageCommands(ctx context.Context, message *tgbotapi.Message) string {

	var message_text string

	switch message.Command() {
	case "start":
		message_text = "Команда /start"
	case "help":
		message_text = "Команды:\n/start\n/help\n"
	}
	return message_text
}
