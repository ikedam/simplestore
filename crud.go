package simplestore

import (
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
)

// Get retrieves a document from firestore
// o must be a pointer to a struct.
// Fill o with the found document.
func (c *Client) Get(ctx context.Context, o any) error {
	doc, err := c.GetDocumentRefSafe(o)
	if err != nil {
		return err
	}
	docsnap, err := doc.Get(ctx)
	if err != nil {
		return err
	}
	return docsnap.DataTo(o)
}

// GetAll retrieves multiple documents from firestore
// os must be a slice of a pointer to a struct.
// Fill os with found documents.
// Returns slice of found objects.
func (c *Client) GetAll(ctx context.Context, os any) (any, error) {
	docList, err := c.GetDocumentRefListSafe(os)
	if err != nil {
		return nil, err
	}
	if len(docList) <= 0 {
		return nil, nil
	}
	validList := make([]*firestore.DocumentRef, 0, len(docList))
	// make dstList as same type of os
	osRef := reflect.ValueOf(os)
	dstList := reflect.MakeSlice(osRef.Type(), 0, len(docList))
	for idx, doc := range docList {
		if doc == nil {
			continue
		}
		validList = append(validList, doc)
		dstList = reflect.Append(dstList, osRef.Index(idx))
	}
	docsnapList, err := c.FirestoreClient.GetAll(ctx, validList)
	if err != nil {
		return nil, err
	}
	for idx, docsnap := range docsnapList {
		if !docsnap.Exists() {
			continue
		}
		err := docsnap.DataTo(dstList.Index(idx).Interface())
		if err != nil {
			return nil, err
		}
	}
	return dstList.Interface(), nil
}

// Create creates a new document in firestore
// o must be a pointer to a struct.
// Generates and sets ID if not set.
func (c *Client) Create(ctx context.Context, o any) (*firestore.WriteResult, error) {
	doc, setID, err := c.prepareSetDocument(o)
	if err != nil {
		return nil, err
	}
	result, err := doc.Create(ctx, o)
	if err != nil {
		return result, err
	}
	setID()
	return result, nil
}

// Set updates a document if exists nor create a new document
// o must be a pointer to a struct.
// Generates and sets ID if not set.
func (c *Client) Set(ctx context.Context, o any, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	doc, setID, err := c.prepareSetDocument(o)
	if err != nil {
		return nil, err
	}
	result, err := doc.Set(ctx, o, opts...)
	if err != nil {
		return result, err
	}
	setID()
	return result, nil
}

// Delete deletes a document
// o must be a pointer to a struct.
func (c *Client) Delete(ctx context.Context, o any, opts ...firestore.Precondition) (*firestore.WriteResult, error) {
	doc, err := c.GetDocumentRefSafe(o)
	if err != nil {
		return nil, err
	}
	return doc.Delete(ctx, opts...)
}
