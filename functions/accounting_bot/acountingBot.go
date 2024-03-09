package accounting_bot

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cloud.google.com/go/firestore"
)

var FIRESTORE_CLIENT *firestore.Client
var FIRESTORE_ERR error
var TG_BOT *tgbotapi.BotAPI
var TG_ERR error
var U userCtxKey = "USER"

func init() {

	tg_bot_token, env_success := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if !env_success {
		log.Panic("Error getting TELEGRAM_BOT_TOKEN environment variable")
	}

	google_project_id, env_success := os.LookupEnv("GOOGLE_PROJECT_ID")
	if !env_success {
		log.Panic("Error getting GOOGLE_PROJECT_ID environment variable")
	}

	google_firestore_db_id, env_success := os.LookupEnv("GOOGLE_FIRESTORE_DB_ID")
	if !env_success {
		log.Panic("Error getting GOOGLE_FIRESTORE_DB_ID environment variable")
	}

	if TG_BOT == nil || TG_ERR != nil {

		TG_BOT, TG_ERR = tgbotapi.NewBotAPI(tg_bot_token)
		if TG_ERR != nil {
			log.Panicf("Error initializing telegram bot API client: \"%s\"", TG_ERR)
		}

		TG_BOT.Debug = true

		log.Printf("Authorized on account %s", TG_BOT.Self.UserName)

	}

	if FIRESTORE_CLIENT == nil || FIRESTORE_ERR != nil {

		FIRESTORE_CLIENT, FIRESTORE_ERR = firestore.NewClientWithDatabase(context.TODO(), google_project_id, google_firestore_db_id)

		if FIRESTORE_ERR != nil {
			log.Panicf("Error initialising firestore client: \"%s\"", FIRESTORE_ERR)
		}
	}

	functions.HTTP("dispatchMessages", dispatchMessages)

}

func dispatchMessages(w http.ResponseWriter, r *http.Request) {

	pRuntime := newRuntime(r.Context(), FIRESTORE_CLIENT, TG_BOT)

	var err error
	var update *tgbotapi.Update

	update, err = pRuntime.tg.HandleUpdate(r)

	if err != nil {
		log.Fatalf("Function TG_BOT.HandleUpdate(r) returned an error: \"%v\"", err)
	} else {
		if update.Message == nil {
			log.Printf("got update not containing message: %v", *update)
		} else { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			pRuntime.user, err = setupUserContext(pRuntime, update.Message.From.ID)
			if err != nil {
				log.Fatalf("Canot find user with TelegramID: \"%d\"", update.Message.From.ID)
			}

			pRuntime.user.updateStatistics(pRuntime)

			var message_text string

			if update.Message.IsCommand() {
				message_text = processMessageCommands(pRuntime, update.Message)
			} else {
				message_text = processMessageText(pRuntime, update.Message)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message_text)
			msg.ReplyToMessageID = update.Message.MessageID

			pRuntime.tg.Send(msg)
		}
	}

}