package accounting_bot

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	functions.HTTP("dispatchMessages", dispatchMessages)
}

func dispatchMessages(w http.ResponseWriter, r *http.Request) {

	var err error
	var bot *tgbotapi.BotAPI
	var update *tgbotapi.Update

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panicf("Error getting TELEGRAM_BOT_TOKEN environment variable: \"%s\"", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	update, err = bot.HandleUpdate(r)

	if err != nil {
		log.Printf("Function bot.HandleUpdate(r) returned an error: \"%v\"", err)
	} else {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var message_text string

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					message_text = "Команда /start"
				case "help":
					message_text = "Команды:\n/start\n/help\n"
				}
			} else {
				message_text = update.Message.Text
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message_text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		} else {
			log.Printf("got update not containing message: %v", update)
		}
	}

}
