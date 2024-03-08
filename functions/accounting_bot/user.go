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

func (u TelegramUser) updateStatistics(ctx context.Context, client *firestore.Client) (*TelegramUser, error) {
	_, err := client.Collection("users").Doc(u.UserId).Update(ctx, []firestore.Update{
		{Path: "message_count", Value: firestore.Increment(1)},
	})

	return &u, err
}

func (u TelegramUser) getMessageCount(ctx context.Context, client *firestore.Client) (int64, error) {

	dsnapUser, err := client.Collection("users").Doc(u.UserId).Get(ctx)
	if err != nil {
		return 0, err
	}

	var user User
	err = dsnapUser.DataTo(&user)

	return user.MessageCount, err
}

func createTelegramUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*TelegramUser, error) {

	var u TelegramUser
	user := User{}

	err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		drefUser, _, err := client.Collection("users").Add(ctx, user)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		u.UserId = drefUser.ID

		_, err = client.Collection("telegram_users").Doc(fmt.Sprintf("%d", telegram_id)).Set(ctx, u)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
		}

		return err
	})

	return &u, err
}

func getTelegramUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*TelegramUser, error) {

	dsnap, err := client.Collection("telegram_users").Doc(fmt.Sprintf("%d", telegram_id)).Get(ctx)

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

func getOrCreateTelegramUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*TelegramUser, error) {
	u, err := getTelegramUser(ctx, telegram_id, client)
	if err != nil {
		u, err = createTelegramUser(ctx, telegram_id, client)
	}
	return u, err
}
