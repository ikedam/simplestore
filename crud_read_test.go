package simplestore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDocument(t *testing.T) {
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
}

func TestGetMultipleDocuments(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
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

func TestReadWithParent(t *testing.T) {
	clearAllDocuments(t, &ParentDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)

	// Create a parent document
	parent := &ParentDocument{
		ID:   "parent1",
		Name: "ParentName",
	}
	_, err = client.Create(ctx, parent)
	assert.NoError(t, err)

	// Create a child document with the parent
	child := &ChildDocument{
		Parent: parent,
		ID:     "child1",
		Name:   "ChildName",
	}
	_, err = client.Create(ctx, child)
	assert.NoError(t, err)

	// Retrieve the child document
	retrievedChild := &ChildDocument{
		Parent: parent,
		ID:     "child1",
	}
	err = client.Get(ctx, retrievedChild)
	assert.NoError(t, err)
	assert.Equal(t, "ChildName", retrievedChild.Name)
	assert.Equal(t, "parent1", retrievedChild.Parent.ID)
	assert.Equal(t, "ParentName", retrievedChild.Parent.Name)
}
