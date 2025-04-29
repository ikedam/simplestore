package simplestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeSafedClient_Get(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)
	typesafed := TypeSafed[MyDocument](client)

	doc := &MyDocument{
		ID: "docid",
	}
	err = typesafed.Get(ctx, doc)
	assert.NoError(t, err)
	assert.Equal(t, "docid", doc.ID)
}

func TestTypeSafedClient_Set(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)
	typesafed := TypeSafed[MyDocument](client)

	doc := &MyDocument{
		ID:   "docid",
		Name: "Alice",
	}
	_, err = typesafed.Set(ctx, doc)
	assert.NoError(t, err)
}

func TestTypeSafedClient_Delete(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)
	typesafed := TypeSafed[MyDocument](client)

	doc := &MyDocument{
		ID: "docid",
	}
	_, err = typesafed.Delete(ctx, doc)
	assert.NoError(t, err)
}
