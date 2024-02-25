package accounting_bot

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type userCtxKey string

func setupUserContext(ctx context.Context, telegram_user_id int64, dbc *firestore.Client) (context.Context, error) {
	u, err := getOrCreateUser(ctx, telegram_user_id, dbc)

	if err != nil {
		log.Printf("Function getOrCreateUser() returned an error: \"%v\"", err)
		return nil, err
	}

	return context.WithValue(ctx, userCtxKey("USER"), u), nil
}
