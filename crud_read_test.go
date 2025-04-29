package simplestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MyDocument struct {
	ID   string
	Name string
}

func TestGetDocument(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Create a test document
	doc := &MyDocument{
		ID:   "docid",
		Name: "Test Name",
	}
	_, err = client.Set(ctx, doc)
	require.NoError(t, err)

	// Retrieve the document
	retrievedDoc := &MyDocument{ID: "docid"}
	err = client.Get(ctx, retrievedDoc)
	require.NoError(t, err)

	// Assert the retrieved document matches the original
	assert.Equal(t, doc.Name, retrievedDoc.Name)
}

func TestGetMultipleDocuments(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Create test documents
	doc1 := &MyDocument{
		ID:   "docid1",
		Name: "Name1",
	}
	doc2 := &MyDocument{
		ID:   "docid2",
		Name: "Name2",
	}
	_, err = client.Set(ctx, doc1)
	require.NoError(t, err)
	_, err = client.Set(ctx, doc2)
	require.NoError(t, err)

	// Retrieve the documents
	retrievedDocs := []*MyDocument{
		{ID: "docid1"},
		{ID: "docid2"},
	}
	_, err = client.GetAll(ctx, retrievedDocs)
	require.NoError(t, err)

	// Assert the retrieved documents match the originals
	assert.Equal(t, doc1.Name, retrievedDocs[0].Name)
	assert.Equal(t, doc2.Name, retrievedDocs[1].Name)
}
