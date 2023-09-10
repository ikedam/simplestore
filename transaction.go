package simplestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func (c *Client) RunTransaction(ctx context.Context, f func(ctx context.Context, client *Client) error, opts ...firestore.TransactionOption) error {
	return c.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		newClient := *c
		newClient.FirestoreTransaction = t
		return f(ctx, &newClient)
	}, opts...)
}
