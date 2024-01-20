package accounting_bot

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

type User struct {
	TelegramId int64 `json:"telegram_id,omitempty" firestore:"telegram_id,omitempty"`
}

func createUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*User, error) {

	var u User = User{TelegramId: telegram_id}

	_, _, err := client.Collection("users").Add(ctx, u)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}

	return &u, err
}

func getUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*User, error) {
	dsnaps, err := client.Collection("users").Where("telegram_id", "==", telegram_id).Documents(ctx).GetAll()

	if err != nil {
		return nil, err
	}

	if len(dsnaps) > 1 {
		return nil, fmt.Errorf("getUser(): multiple users found with telegram_id==%d", telegram_id)
	}

	if len(dsnaps) < 1 {
		return nil, fmt.Errorf("getUser(): user not found with telegram_id==%d", telegram_id)
	}

	var u User
	err = dsnaps[0].DataTo(&u)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Document data: %#v\n", u)
	return &u, nil
}

func getOrCreateUser(ctx context.Context, telegram_id int64, client *firestore.Client) (*User, error) {
	u, err := getUser(ctx, telegram_id, client)
	if err != nil {
		u, err = createUser(ctx, telegram_id, client)
	}
	return u, err
}
