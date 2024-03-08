package accounting_bot

import (
	"log"
)

type userCtxKey string

func setupUserContext(runtime *botRuntime, telegram_user_id int64) (*TelegramUser, error) {
	u, err := getOrCreateTelegramUser(runtime, telegram_user_id)

	if err != nil {
		log.Printf("Function getOrCreateUser() returned an error: \"%v\"", err)
		return nil, err
	}

	return u, nil
}
