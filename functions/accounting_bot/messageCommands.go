package accounting_bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func processMessageCommands(runtime *botRuntime, message *tgbotapi.Message) string {

	var message_text string

	switch message.Command() {
	case "start":
		message_text = "Команда /start"
	case "help":
		message_text = "Команды:\n/start\n/help\n"

		for k, v := range runtime.r.Header {
			message_text += fmt.Sprintf("Header: %s, value: %v\n", k, v)
		}

		message_text += fmt.Sprintf("\nYour Client IP: %s\n", runtime.r.RemoteAddr)

	}
	return message_text
}
