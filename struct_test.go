// filepath: /workspaces/simplestore/struct_test.go
package simplestore

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CustomIDDocument implements the IDer interface to customize document ID handling
type CustomIDDocument struct {
	MyID string
	Name string
}

// GetDocumentID generates a document ID with a hash prefix based on MyID
func (d *CustomIDDocument) GetDocumentID() string {
	if d.MyID == "" {
		return ""
	}
	// Simple hash prefix for demonstration
	return fmt.Sprintf("hash_%s", d.MyID)
}

// SetDocumentID extracts the original MyID from the document ID
func (d *CustomIDDocument) SetDocumentID(id string) {
	if id == "" {
		d.MyID = ""
		return
	}

	// Extract the original ID by removing the hash_ prefix
	if strings.HasPrefix(id, "hash_") {
		d.MyID = id[5:] // Remove "hash_" prefix
	} else {
		// Unexpected format
		d.MyID = id
	}
}

func TestCustomIDDocumentInterface(t *testing.T) {
	// Test that CustomIDDocument implements the IDer interface
	var doc interface{} = &CustomIDDocument{}
	_, ok := doc.(IDer)
	assert.True(t, ok, "CustomIDDocument should implement IDer interface")

	// Test GetDocumentID
	testDoc := &CustomIDDocument{MyID: "test123"}
	assert.Equal(t, "hash_test123", testDoc.GetDocumentID())

	// Test SetDocumentID
	testDoc = &CustomIDDocument{}
	testDoc.SetDocumentID("hash_test456")
	assert.Equal(t, "test456", testDoc.MyID)
}

func TestCreateCustomIDDocument(t *testing.T) {
	clearAllDocuments(t, &CustomIDDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Create a new document with custom ID
	doc := &CustomIDDocument{
		MyID: "custom123",
		Name: "Test Document",
	}

	_, err = client.Create(ctx, doc)
	assert.NoError(t, err)
	assert.Equal(t, "custom123", doc.MyID, "MyID should remain unchanged")

	// Verify document ID in Firestore would be "hash_custom123"
	expectedID := "hash_custom123"

	// Retrieve the document using the expected ID format
	retrievedDoc := &CustomIDDocument{}
	retrievedDoc.SetDocumentID(expectedID)
	err = client.Get(ctx, retrievedDoc)
	require.NoError(t, err)

	// Assert the retrieved document matches the original
	assert.Equal(t, doc.MyID, retrievedDoc.MyID)
	assert.Equal(t, doc.Name, retrievedDoc.Name)
}

func TestGetCustomIDDocument(t *testing.T) {
	clearAllDocuments(t, &CustomIDDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// Create a document to retrieve later
	doc := &CustomIDDocument{
		MyID: "get456",
		Name: "Get Test",
	}

	_, err = client.Create(ctx, doc)
	require.NoError(t, err)

	// Create a new instance to retrieve the document
	retrieveDoc := &CustomIDDocument{
		MyID: "get456",
	}

	// Retrieve the document
	err = client.Get(ctx, retrieveDoc)
	assert.NoError(t, err)

	// Check that both MyID and Name were properly retrieved
	assert.Equal(t, "get456", retrieveDoc.MyID)
	assert.Equal(t, "Get Test", retrieveDoc.Name)
}

func TestSetCustomIDDocument(t *testing.T) {
	clearAllDocuments(t, &CustomIDDocument{})
	ctx := context.Background()
	client, err := New(ctx)
	require.NoError(t, err)

	// First create a document
	doc := &CustomIDDocument{
		MyID: "set789",
		Name: "Original Name",
	}

	_, err = client.Create(ctx, doc)
	require.NoError(t, err)

	// Update the document using Set
	updatedDoc := &CustomIDDocument{
		MyID: "set789",
		Name: "Updated Name",
	}

	_, err = client.Set(ctx, updatedDoc)
	assert.NoError(t, err)

	// Retrieve the document to verify the update
	retrievedDoc := &CustomIDDocument{
		MyID: "set789",
	}

	err = client.Get(ctx, retrievedDoc)
	require.NoError(t, err)

	// Check that the name was updated
	assert.Equal(t, "Updated Name", retrievedDoc.Name)
}
