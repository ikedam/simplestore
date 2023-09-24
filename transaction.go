package simplestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

func (c *Client) RunTransaction(ctx context.Context, f func(ctx context.Context, client *Client) error, opts ...firestore.TransactionOption) error {
	newClient := *c
	err := c.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		newClient.FirestoreTransaction = t
		return f(ctx, &newClient)
	}, opts...)
	if err != nil {
		for _, callback := range newClient.transactionFailureCallbacks {
			callback()
		}
	}
	return err
}
