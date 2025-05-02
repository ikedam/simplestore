package simplestore

// IDer is an interface that defines a method for getting and setting the firestore document ID of a struct.
type IDer interface {
	GetDocumentID() string
	SetDocumentID(id string)
}
