package accounting_bot

import (
	"context"
	"fmt"
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
var PP_STRIPE_TOKEN string

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

	PP_STRIPE_TOKEN, env_success = os.LookupEnv("PP_STRIPE_TOKEN")
	if !env_success {
		log.Panic("Error getting PP_STRIPE_TOKEN environment variable")
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
		if update.Message != nil { // If we got a message
			log.Printf("Message from [%s] with text: \"%s\"", update.Message.From.UserName, update.Message.Text)

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
		} else if update.InlineQuery != nil {
			log.Printf("[%s] %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)

			pRuntime.user, err = setupUserContext(pRuntime, update.InlineQuery.From.ID)
			if err != nil {
				log.Fatalf("Canot find user with TelegramID: \"%d\"", update.InlineQuery.From.ID)
			}

			pRuntime.user.updateStatistics(pRuntime)

			invc := tgbotapi.InputInvoiceMessageContent{
				ProviderToken: PP_STRIPE_TOKEN,
				Title:         "Here is the prices",
				Currency:      "USD",
				Prices: []tgbotapi.LabeledPrice{
					{Label: "Per hour", Amount: 7500},
					{Label: "Per day", Amount: 100000},
					{Label: "Per project", Amount: 1000000},
				},
			}

			res := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID,
				"Prices Article Title",
				fmt.Sprintf("Here are the prices for your request: \"%s\"", update.InlineQuery.Query))

			res.InputMessageContent = invc

			inlineConf := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{res},
			}

			if _, err := pRuntime.tg.Request(inlineConf); err != nil {
				log.Printf("Failed to send inline query response: %s", err)
			}

		} else {
			log.Printf("got update not containing message nor inline query: %v", *update)
		}
	}

}
