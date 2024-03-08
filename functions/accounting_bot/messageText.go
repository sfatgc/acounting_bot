package accounting_bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processMessageText(runtime *botRuntime, message *tgbotapi.Message) string {
	u := runtime.user
	message_count, err := u.getMessageCount(runtime)
	if err != nil {
		return fmt.Sprintf("NaN (%s)", err)
	}
	return fmt.Sprintf("You wrote:\n\t```\n%s\n```\nYou sent %d messages", message.Text, message_count)
}
