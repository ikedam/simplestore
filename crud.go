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
	var docsnap *firestore.DocumentSnapshot
	if c.FirestoreTransaction == nil {
		docsnap, err = doc.Get(ctx)
	} else {
		docsnap, err = c.FirestoreTransaction.Get(doc)
	}
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
	var docsnapList []*firestore.DocumentSnapshot
	if c.FirestoreTransaction == nil {
		docsnapList, err = c.FirestoreClient.GetAll(ctx, validList)
	} else {
		docsnapList, err = c.FirestoreTransaction.GetAll(validList)
	}
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
// WriteResult will be alwasys `nil` while transaction.
func (c *Client) Create(ctx context.Context, o any) (*firestore.WriteResult, error) {
	accessor, err := newAccessor(reflect.TypeOf(o), c.tableMaps)
	if err != nil {
		return nil, err
	}
	if accessor.readOnly {
		return nil, NewProgrammingErrorf("cannot create document in readonly collection: %s", accessor.collectionName)
	}

	doc, resetID, err := c.prepareSetDocument(o)
	if err != nil {
		return nil, err
	}
	var result *firestore.WriteResult
	if c.FirestoreTransaction == nil {
		result, err = doc.Create(ctx, o)
		if err != nil {
			resetID()
			return result, err
		}
	} else {
		err = c.FirestoreTransaction.Create(doc, o)
		if err != nil {
			return result, err
		}
		c.transactionFailureCallbacks = append(c.transactionFailureCallbacks, resetID)
	}
	return result, nil
}

// Set updates a document if exists nor create a new document
// o must be a pointer to a struct.
// Generates and sets ID if not set.
// WriteResult will be alwasys `nil` while transaction.
func (c *Client) Set(ctx context.Context, o any, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	accessor, err := newAccessor(reflect.TypeOf(o), c.tableMaps)
	if err != nil {
		return nil, err
	}
	if accessor.readOnly {
		return nil, NewProgrammingErrorf("cannot set document in readonly collection: %s", accessor.collectionName)
	}

	doc, resetID, err := c.prepareSetDocument(o)
	if err != nil {
		return nil, err
	}
	var result *firestore.WriteResult
	if c.FirestoreTransaction == nil {
		result, err = doc.Set(ctx, o, opts...)
		if err != nil {
			resetID()
			return result, err
		}
	} else {
		err = c.FirestoreTransaction.Set(doc, o)
		if err != nil {
			return result, err
		}
		c.transactionFailureCallbacks = append(c.transactionFailureCallbacks, resetID)
	}
	return result, nil
}

// Delete deletes a document
// o must be a pointer to a struct.
// WriteResult will be alwasys `nil` while transaction.
func (c *Client) Delete(ctx context.Context, o any, opts ...firestore.Precondition) (*firestore.WriteResult, error) {
	accessor, err := newAccessor(reflect.TypeOf(o), c.tableMaps)
	if err != nil {
		return nil, err
	}
	if accessor.readOnly {
		return nil, NewProgrammingErrorf("cannot delete document in readonly collection: %s", accessor.collectionName)
	}

	doc, err := c.GetDocumentRefSafe(o)
	if err != nil {
		return nil, err
	}
	if c.FirestoreTransaction == nil {
		return doc.Delete(ctx, opts...)
	} else {
		return nil, c.FirestoreTransaction.Delete(doc, opts...)
	}
}
