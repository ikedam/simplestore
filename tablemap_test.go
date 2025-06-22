package simplestore

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTableMapDocument is a test struct for table mapping
type TestTableMapDocument struct {
	ID   string
	Name string
}

// TestReadOnlyDocument is a test struct for readonly collection
type TestReadOnlyDocument struct {
	ID   string
	Name string
}

func TestAddTableMaps(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Test adding table maps
	tableMap := map[string]string{
		"TestTableMapDocument": "custom_collection",
		"TestReadOnlyDocument": "readonly_collection",
	}

	client.AddTableMaps(tableMap)

	// Verify table maps were added
	assert.Equal(t, "custom_collection", client.tableMaps["TestTableMapDocument"].CollectionName)
	assert.False(t, client.tableMaps["TestTableMapDocument"].ReadOnly)
	assert.Equal(t, "readonly_collection", client.tableMaps["TestReadOnlyDocument"].CollectionName)
	assert.False(t, client.tableMaps["TestReadOnlyDocument"].ReadOnly)
}

func TestAddReadonlyTableMaps(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Test adding readonly table maps
	tableMap := map[string]string{
		"TestReadOnlyDocument": "readonly_collection",
	}

	client.AddReadonlyTableMaps(tableMap)

	// Verify readonly table maps were added
	assert.Equal(t, "readonly_collection", client.tableMaps["TestReadOnlyDocument"].CollectionName)
	assert.True(t, client.tableMaps["TestReadOnlyDocument"].ReadOnly)
}

func TestTableMapCollectionName(t *testing.T) {
	clearAllDocuments(t, &TestTableMapDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Add table mapping
	tableMap := map[string]string{
		"TestTableMapDocument": "custom_collection",
	}
	client.AddTableMaps(tableMap)

	// Create a document
	doc := &TestTableMapDocument{
		Name: "Test Document",
	}
	_, err = client.Create(ctx, doc)
	require.NoError(t, err)

	// Verify the document was created in the custom collection
	// The document should be accessible via the custom collection name
	accessor, err := newAccessor(reflect.TypeOf(doc), client.tableMaps)
	require.NoError(t, err)
	assert.Equal(t, "custom_collection", accessor.collectionName)
}

func TestReadOnlyCollection(t *testing.T) {
	clearAllDocuments(t, &TestReadOnlyDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Add readonly table mapping
	tableMap := map[string]string{
		"TestReadOnlyDocument": "readonly_collection",
	}
	client.AddReadonlyTableMaps(tableMap)

	// Test that Create fails for readonly collection
	doc := &TestReadOnlyDocument{
		Name: "Test Document",
	}
	_, err = client.Create(ctx, doc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create document in readonly collection")

	// Test that Set fails for readonly collection
	doc2 := &TestReadOnlyDocument{
		ID:   "test-id",
		Name: "Test Document",
	}
	_, err = client.Set(ctx, doc2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot set document in readonly collection")

	// Test that Delete fails for readonly collection
	_, err = client.Delete(ctx, doc2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete document in readonly collection")
}

func TestReadOnlyCollectionWithAccessor(t *testing.T) {
	clearAllDocuments(t, &MyDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Add readonly table mapping
	tableMap := map[string]string{
		"TestReadOnlyDocument": "readonly_collection",
	}
	client.AddReadonlyTableMaps(tableMap)

	// Test that accessor correctly sets readonly flag
	doc := &TestReadOnlyDocument{
		ID:   "test-id",
		Name: "Test Document",
	}
	accessor, err := newAccessor(reflect.TypeOf(doc), client.tableMaps)
	require.NoError(t, err)
	assert.Equal(t, "readonly_collection", accessor.collectionName)
	assert.True(t, accessor.readOnly)
}

func TestTableMapWithQuery(t *testing.T) {
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)
	require.NoError(t, DeleteCollection(ctx, client.FirestoreClient, "custom_collection", 100))

	// Add table mapping
	tableMap := map[string]string{
		"TestTableMapDocument": "custom_collection",
	}
	client.AddTableMaps(tableMap)

	// Create a document
	doc := &TestTableMapDocument{
		Name: "Test Document",
	}
	_, err = client.Create(ctx, doc)
	require.NoError(t, err)

	// Test query with table mapping
	var docs []*TestTableMapDocument
	q := client.Query(&docs)
	err = q.GetAll(ctx)
	require.NoError(t, err)

	// Verify the query worked with the custom collection name
	assert.Len(t, docs, 1)
	assert.Equal(t, "Test Document", docs[0].Name)
}
