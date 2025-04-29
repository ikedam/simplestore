package simplestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateDocument(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	doc := &MyDocument{
		Name: "Alice",
	}

	_, err = client.Create(ctx, doc)
	assert.NoError(t, err)
	assert.NotEmpty(t, doc.ID) // Ensure ID is generated if not set

	// Retrieve the document
	retrievedDoc := &MyDocument{ID: doc.ID}
	err = client.Get(ctx, retrievedDoc)
	require.NoError(t, err)

	// Assert the retrieved document matches the original
	assert.Equal(t, doc.Name, retrievedDoc.Name)
}

func TestSetDocument(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	doc1 := &MyDocument{
		Name: "Alice",
	}

	_, err = client.Create(ctx, doc1)
	assert.NoError(t, err)

	doc2 := &MyDocument{
		ID:   doc1.ID,
		Name: "Bob",
	}

	_, err = client.Set(ctx, doc2)
	assert.NoError(t, err)

	// Retrieve the document
	retrievedDoc := &MyDocument{ID: doc2.ID}
	err = client.Get(ctx, retrievedDoc)
	require.NoError(t, err)

	// Assert the retrieved document matches the original
	assert.Equal(t, doc2.Name, retrievedDoc.Name)
}

func TestDeleteDocument(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
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

	_, err = client.Delete(ctx, doc)
	assert.NoError(t, err)

	err = client.Get(ctx, retrievedDoc)
	// grpc の NotFound 応答
	require.Equal(t, status.Code(err), codes.NotFound)
}
