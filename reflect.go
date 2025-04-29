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
	if pt.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
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
	return collection.Doc(docID), false, nil
}

func (a *accessor) setID(pv reflect.Value, id string) {
	v := pv.Elem()
	v.FieldByName(IDFieldName).SetString(id)
}

// GetDocumentRefSafe returns document ref of the object
// o must be a pointer to a struct.
// Returns nil if object is a nil.
// Error if ID is not set.
func (c *Client) GetDocumentRefSafe(o any) (*firestore.DocumentRef, error) {
	accessor, err := newAccessor(reflect.TypeOf(o))
	if err != nil {
		return nil, err
	}
	doc, _, err := accessor.getDocumentRef(c, reflect.ValueOf(o), false)
	return doc, err
}

func nop() {}

func (c *Client) prepareSetDocument(o any) (*firestore.DocumentRef, func(), error) {
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
	accessor.setID(pv, doc.ID)
	return doc, func() {
		// reset
		accessor.setID(pv, "")
	}, nil
}

// GetDocumentRef returns document ref of the object
// o must be a pointer to a struct.
// Returns nil if object is a nil.
// Panic if inappropriate value passed.
func (c *Client) GetDocumentRef(o any) *firestore.DocumentRef {
	doc, err := c.GetDocumentRefSafe(o)
	if err != nil {
		panic(err)
	}
	return doc
}

// GetDocumentRefListSafe returns a list of document refs of the objects
// os must be a slice of a pointer to a struct.
// Be aware that return value may contain `nil`.
func (c *Client) GetDocumentRefListSafe(os any) ([]*firestore.DocumentRef, error) {
	osT := reflect.TypeOf(os)
	if osT.Kind() != reflect.Slice {
		return nil, NewProgrammingError("expects a slice of pointers to structs")
	}
	poT := osT.Elem()
	if poT.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("expects a slice of pointers to structs")
	}
	elemType := poT.Elem()
	if elemType.Kind() == reflect.Struct {
		// all elements are same type
		return c.getDocumentRefListSafeWithSameType(os)
	}
	var docList []*firestore.DocumentRef
	for _, o := range os.([]any) {
		doc, err := c.GetDocumentRefSafe(o)
		if err != nil {
			return nil, err
		}
		docList = append(docList, doc)
	}
	return docList, nil
}

func (c *Client) getDocumentRefListSafeWithSameType(os any) ([]*firestore.DocumentRef, error) {
	sliceRef := reflect.ValueOf(os)
	accessor, err := newAccessor(reflect.TypeOf(os).Elem())
	if err != nil {
		return nil, err
	}

	var docList []*firestore.DocumentRef
	for i := 0; i < sliceRef.Len(); i++ {
		doc, _, err := accessor.getDocumentRef(c, sliceRef.Index(i), false)
		if err != nil {
			return nil, err
		}
		docList = append(docList, doc)
	}
	return docList, nil
}

type targetBuilder struct {
	parent         any
	target         any
	elementType    reflect.Type
	collectionName string
}

func newTargetBuilder(pos any) (*targetBuilder, error) {
	posT := reflect.TypeOf(pos)
	if posT.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("expects a pointer to a slice of pointers to structs")
	}
	osT := posT.Elem()
	if osT.Kind() != reflect.Slice {
		return nil, NewProgrammingError("expects a pointer to a slice of pointers to structs")
	}
	poT := osT.Elem()
	if poT.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	t := poT.Elem()
	if t.Kind() != reflect.Struct {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	return &targetBuilder{
		target:         pos,
		elementType:    t,
		collectionName: t.Name(),
	}, nil
}

func newTargetBuilderWithParent(parent any, pos any) (*targetBuilder, error) {
	posT := reflect.TypeOf(pos)
	if posT.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("expects a pointer to a slice of pointers to structs")
	}
	osT := posT.Elem()
	if osT.Kind() != reflect.Slice {
		return nil, NewProgrammingError("expects a pointer to a slice of pointers to structs")
	}
	poT := osT.Elem()
	if poT.Kind() != reflect.Pointer {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	t := poT.Elem()
	if t.Kind() != reflect.Struct {
		return nil, NewProgrammingError("value must be a pointer of a struct")
	}
	parentT := reflect.TypeOf(parent)
	parentF, ok := t.FieldByName(ParentFieldName)
	if !ok {
		return nil, NewProgrammingErrorf("value must have "+ParentFieldName+" field with type %s.%s", parentT.PkgPath(), parentT.Name())
	}
	if !parentT.AssignableTo(parentF.Type) {
		return nil, NewProgrammingErrorf(ParentFieldName+" field must be %s.%s", parentT.PkgPath(), parentT.Name())
	}
	return &targetBuilder{
		target:         pos,
		parent:         parent,
		elementType:    t,
		collectionName: t.Name(),
	}, nil
}

func (t *targetBuilder) createElement() any {
	pv := reflect.New(t.elementType)
	if t.parent != nil {
		pv.Elem().FieldByName(ParentFieldName).Set(reflect.ValueOf(t.parent))
	}
	return pv.Interface()
}

func (t *targetBuilder) append(o any) {
	refSlice := reflect.ValueOf(t.target).Elem()
	refSlice.Set(reflect.Append(
		refSlice,
		reflect.ValueOf(o),
	))
}

func (c *Client) getCollectionRef(pos any) (*firestore.CollectionRef, *targetBuilder, error) {
	tb, err := newTargetBuilder(pos)
	if err != nil {
		return nil, nil, err
	}
	return c.FirestoreClient.Collection(tb.collectionName), tb, nil
}

func (c *Client) getCollectionGroupRef(pos any) (*firestore.CollectionGroupRef, *targetBuilder, error) {
	tb, err := newTargetBuilder(pos)
	if err != nil {
		return nil, nil, err
	}
	return c.FirestoreClient.CollectionGroup(tb.collectionName), tb, nil
}

func (c *Client) getNestedCollectionRef(parent any, pos any) (*firestore.CollectionRef, *targetBuilder, error) {
	tb, err := newTargetBuilderWithParent(parent, pos)
	if err != nil {
		return nil, nil, err
	}

	doc, err := c.GetDocumentRefSafe(parent)
	if err != nil {
		return nil, nil, err
	}
	return doc.Collection(tb.collectionName), tb, nil
}
