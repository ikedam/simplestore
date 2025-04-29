package simplestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTransaction(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)

	doc := &MyDocument{
		ID:   "docid",
		Name: "Alice",
	}

	// Create the document initially
	_, err = client.Create(ctx, doc)
	assert.NoError(t, err)

	err = client.RunTransaction(ctx, func(ctx context.Context, txClient *Client) error {
		// Retrieve the document within the transaction
		err := txClient.Get(ctx, doc)
		if err != nil {
			return err
		}

		// Update the document's name
		doc.Name = "Bob"
		_, err = txClient.Set(ctx, doc)
		return err
	})
	assert.NoError(t, err)

	// Verify the document was updated
	updatedDoc := &MyDocument{ID: "docid"}
	err = client.Get(ctx, updatedDoc)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", updatedDoc.Name)
}
