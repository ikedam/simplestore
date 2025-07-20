package simplestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
)

func TestQuerySimpleDocuments(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)

	// Prepare test data
	doc1 := &MyDocument{ID: "docid1", Name: "Alice"}
	doc2 := &MyDocument{ID: "docid2", Name: "Bob"}
	_, err = client.Set(ctx, doc1)
	assert.NoError(t, err)
	_, err = client.Set(ctx, doc2)
	assert.NoError(t, err)

	// Query documents
	var docs []*MyDocument
	q := client.Query(&docs).OrderBy("Name", firestore.Asc)
	err = q.GetAll(ctx)
	assert.NoError(t, err)

	// Validate results
	assert.Len(t, docs, 2)
	assert.Equal(t, "Alice", docs[0].Name)
	assert.Equal(t, "Bob", docs[1].Name)
}

func TestQueryNestedDocuments(t *testing.T) {
	clearAllDocuments(t, &ParentDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)

	// Prepare test data
	parent := &ParentDocument{ID: "parent1", Name: "Parent"}
	child1 := &ChildDocument{Parent: parent, ID: "child1", Name: "Child1"}
	child2 := &ChildDocument{Parent: parent, ID: "child2", Name: "Child2"}
	_, err = client.Set(ctx, parent)
	assert.NoError(t, err)
	_, err = client.Set(ctx, child1)
	assert.NoError(t, err)
	_, err = client.Set(ctx, child2)
	assert.NoError(t, err)

	// Query nested documents
	var childDocs []*ChildDocument
	q := client.QueryNested(parent, &childDocs)
	err = q.GetAll(ctx)
	assert.NoError(t, err)

	// Validate results
	assert.Len(t, childDocs, 2)
	assert.Equal(t, "Child1", childDocs[0].Name)
	assert.Equal(t, "Child2", childDocs[1].Name)
}

func TestCount(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	assert.NoError(t, err)

	clearAllDocuments(t, &MyDocument{})

	// Query documents
	var docs []*MyDocument
	q := client.Query(&docs)
	count, err := q.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Prepare test data
	doc1 := &MyDocument{ID: "docid1", Name: "Alice"}
	doc2 := &MyDocument{ID: "docid2", Name: "Bob"}
	_, err = client.Set(ctx, doc1)
	assert.NoError(t, err)
	_, err = client.Set(ctx, doc2)
	assert.NoError(t, err)

	// Query documents
	count, err = q.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Query documents
	count, err = q.Where("Name", "==", "Alice").Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
