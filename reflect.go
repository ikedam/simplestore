package simplestore

import (
	"reflect"

	"cloud.google.com/go/firestore"
)

const (
	IDFieldName     = "ID"
	ParentFieldName = "Parent"
)

type accessor struct {
	parentAccessor *accessor
	t              reflect.Type
	collectionName string
}

func newAccessor(pt reflect.Type) (*accessor, error) {
	// pt must be a pointer.
	t := pt.Elem()
	a := &accessor{
		t: t,
	}
	if t.Kind() != reflect.Struct {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	idF, ok := t.FieldByName(IDFieldName)
	if !ok {
		return nil, NewProgrammingErrorf(IDFieldName+" field doesn't exist: %s.%s", t.PkgPath(), t.Name())
	}
	if idF.Type.Kind() != reflect.String {
		return nil, NewProgrammingErrorf(IDFieldName+" field must be a string: %s.%s", t.PkgPath(), t.Name())
	}
	a.collectionName = t.Name()

	parentF, ok := t.FieldByName("Parent")
	if ok {
		parentT := parentF.Type
		if parentT.Kind() != reflect.Pointer {
			return nil, NewProgrammingError("Parent must be a pointer of a struct")
		}
		var err error
		a.parentAccessor, err = newAccessor(parentT)
		if err != nil {
			return nil, NewProgrammingErrorf("invalid parent in %s.%s: %v", t.PkgPath(), t.Name(), err.Error())
		}
	}
	return a, nil
}

func (a *accessor) getDocumentRef(c *Client, pv reflect.Value, mightNew bool) (*firestore.DocumentRef, bool, error) {
	if pv.IsNil() {
		return nil, false, nil
	}
	v := pv.Elem()
	docID := v.FieldByName(IDFieldName).String()
	var collection *firestore.CollectionRef

	if a.parentAccessor != nil {
		pparent := v.FieldByName(ParentFieldName)
		parentDoc, _, err := a.parentAccessor.getDocumentRef(c, pparent, false)
		if err != nil {
			return nil, false, NewProgrammingErrorf("invalid parent in %s.%s: %v", a.t.PkgPath(), a.t.Name(), err.Error())
		}
		if parentDoc != nil {
			collection = parentDoc.Collection(a.collectionName)
		}
	}

	if collection == nil {
		collection = c.FirestoreClient.Collection(a.collectionName)
	}
	if docID == "" {
		if mightNew {
			return collection.NewDoc(), true, nil
		} else {
			return nil, false, NewProgrammingError("ID is not set")
		}
	}
	return collection.Doc(docID), true, nil
}

func (a *accessor) setID(pv reflect.Value, id string) {
	v := pv.Elem()
	v.FieldByName(IDFieldName).SetString(id)
}

// GetDocumentRefSafe returns document ref of the object
// Returns nil if object is a nil.
// Error if ID is not set.
func (c *Client) GetDocumentRefSafe(o *any) (*firestore.DocumentRef, error) {
	if o == nil {
		// fast return
		return nil, nil
	}
	accessor, err := newAccessor(reflect.TypeOf(o))
	if err != nil {
		return nil, err
	}
	doc, _, err := accessor.getDocumentRef(c, reflect.ValueOf(o), false)
	return doc, err
}

func nop() {}

func (c *Client) prepareSetDocument(o *any) (*firestore.DocumentRef, func(), error) {
	accessor, err := newAccessor(reflect.TypeOf(o))
	if err != nil {
		return nil, nil, err
	}
	pv := reflect.ValueOf(o)
	doc, isNew, err := accessor.getDocumentRef(c, pv, true)
	if err != nil {
		return nil, nil, err
	}
	if doc == nil {
		return nil, nil, NewProgrammingError("object is nil")
	}
	if !isNew {
		return doc, nop, nil
	}
	return doc, func() {
		accessor.setID(pv, doc.ID)
	}, nil
}

// GetDocumentRef returns document ref of the object
// Returns nil if object is a nil.
// Panic if inappropriate value passed.
func (c *Client) GetDocumentRef(o *any) *firestore.DocumentRef {
	doc, err := c.GetDocumentRefSafe(o)
	if err != nil {
		panic(err)
	}
	return doc
}

// GetDocumentRefListSafe returns a list of document refs of the objects
// Be aware that return value may contain `nil`.
func (c *Client) GetDocumentRefListSafe(os []*any) ([]*firestore.DocumentRef, error) {
	elemType := reflect.TypeOf(os).Elem().Elem() // type of *(os[idx])
	if elemType.Kind() == reflect.Struct {
		// all elements are same type
		return c.getDocumentRefListSafeWithSameType(os)
	}
	var docList []*firestore.DocumentRef
	for _, o := range os {
		doc, err := c.GetDocumentRefSafe(o)
		if err != nil {
			return nil, err
		}
		docList = append(docList, doc)
	}
	return docList, nil
}

func (c *Client) getDocumentRefListSafeWithSameType(os []*any) ([]*firestore.DocumentRef, error) {
	accessor, err := newAccessor(reflect.TypeOf(os).Elem())
	if err != nil {
		return nil, err
	}

	var docList []*firestore.DocumentRef
	for _, o := range os {
		doc, _, err := accessor.getDocumentRef(c, reflect.ValueOf(o), false)
		if err != nil {
			return nil, err
		}
		docList = append(docList, doc)
	}
	return docList, nil
}

func (c *Client) getCollectionRef(os []*any) (*firestore.CollectionRef, error) {
	t := reflect.TypeOf(os).Elem().Elem()
	if t.Kind() != reflect.Struct {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	return c.FirestoreClient.Collection(t.Name()), nil
}
