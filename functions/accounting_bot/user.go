package accounting_bot

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

type User struct {
	MessageCount int64 `json:"message_count,omitempty" firestore:"message_count,omitempty"`
}

type TelegramUser struct {
	UserId string `json:"user_id,omitempty" firestore:"user_id,omitempty"`
}

func (u TelegramUser) updateStatistics(runtime *botRuntime) (*TelegramUser, error) {
	_, err := runtime.db.Collection("users").Doc(u.UserId).Update(runtime.rCtx, []firestore.Update{
		{Path: "message_count", Value: firestore.Increment(1)},
	})

	return &u, err
}

func (u TelegramUser) getMessageCount(runtime *botRuntime) (int64, error) {

	dsnapUser, err := runtime.db.Collection("users").Doc(u.UserId).Get(runtime.rCtx)
	if err != nil {
		return 0, err
	}

	var user User
	err = dsnapUser.DataTo(&user)

	return user.MessageCount, err
}

func createTelegramUser(runtime *botRuntime, telegram_id int64) (*TelegramUser, error) {

	var u TelegramUser
	user := User{}

	err := runtime.db.RunTransaction(runtime.rCtx, func(ctx context.Context, tx *firestore.Transaction) error {

		drefUser, _, err := runtime.db.Collection("users").Add(ctx, user)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		u.UserId = drefUser.ID

		_, err = runtime.db.Collection("telegram_users").Doc(fmt.Sprintf("%d", telegram_id)).Set(ctx, u)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return err
	})

	return &u, err
}

func getTelegramUser(runtime *botRuntime, telegram_id int64) (*TelegramUser, error) {

	dsnap, err := runtime.db.Collection("telegram_users").Doc(fmt.Sprintf("%d", telegram_id)).Get(runtime.rCtx)

	if err != nil {
		return nil, err
	}

	var u TelegramUser
	err = dsnap.DataTo(&u)

	if err != nil {
		return nil, err
	}

	log.Printf("Document data: %#v\n", u)
	return &u, nil
}

func getOrCreateTelegramUser(runtime *botRuntime, telegram_id int64) (*TelegramUser, error) {
	u, err := getTelegramUser(runtime, telegram_id)
	if err != nil {
		u, err = createTelegramUser(runtime, telegram_id)
	}
	return u, err
}
