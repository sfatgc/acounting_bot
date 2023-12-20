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

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	update, err := bot.HandleUpdate(r)

	if err != nil && update.Message != nil { // If we got a message
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

}
