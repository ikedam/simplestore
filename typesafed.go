package simplestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

// TypeSafedClient is a type-restricting wrapper of Client
type TypeSafedClient[T any, P any] struct {
	untyped *Client
}

// TypeSafed returns a TypeSafeClient of Client
func TypeSafed[T any](c *Client) TypeSafedClient[T, any] {
	return TypeSafedClient[T, any]{
		untyped: c,
	}
}

// TypeSafedWithParent returns a TypeSafeClient of Client
func TypeSafedWithParent[T any, P any](c *Client) TypeSafedClient[T, P] {
	return TypeSafedClient[T, P]{
		untyped: c,
	}
}

// GetDocumentRef returns document ref of the object
// Returns nil if object is a nil.
func (c *TypeSafedClient[T, P]) GetDocumentRef(o *T) *firestore.DocumentRef {
	return c.untyped.GetDocumentRef(o)
}

// Get retrieves a document from firestore
// Fill o with the found document.
func (c *TypeSafedClient[T, P]) Get(ctx context.Context, o *T) error {
	return c.untyped.Get(ctx, o)
}

// GetAll retrieves multiple documents from firestore
// Fill os with found documents.
// Returns slice of found objects.
func (c *TypeSafedClient[T, P]) GetAll(ctx context.Context, os []*T) ([]*T, error) {
	res, err := c.untyped.GetAll(ctx, os)
	if res == nil {
		return nil, err
	}
	return res.([]*T), err
}

// Create creates a new document in firestore
// Generates and sets ID if not set.
func (c *TypeSafedClient[T, P]) Create(ctx context.Context, o *T) (*firestore.WriteResult, error) {
	return c.untyped.Create(ctx, o)
}

// Set updates a document if exists nor create a new document
// Generates and sets ID if not set.
func (c *TypeSafedClient[T, P]) Set(ctx context.Context, o *T, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	return c.untyped.Set(ctx, o, opts...)
}

// Delete deletes a document
// o must be a pointer to a struct.
func (c *TypeSafedClient[T, P]) Delete(ctx context.Context, o *T, opts ...firestore.Precondition) (*firestore.WriteResult, error) {
	return c.untyped.Delete(ctx, o, opts...)
}

type TypeSafedQuery[T any] struct {
	untyped *Query
}

// Query starts a new query for target
// target is also used as destination of `GetAll()`.
// Panic if inappropriate target is specified.
func (c *TypeSafedClient[T, P]) Query(target *[]*T) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: c.untyped.Query(target),
	}
}

// QueryGroupSafe starts a new query for target as collection group
// target is also used as destination of `GetAll()`.
func (c *TypeSafedClient[T, P]) QueryGroup(target *[]*T) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: c.untyped.QueryGroup(target),
	}
}

// QueryNested starts a new query for target under the document specified by `parent`
// parent must be a pointer to a struct.
// target is also used as destination of `GetAll()`.
func (c *TypeSafedClient[T, P]) QueryNested(parent *P, target *[]*T) (*TypeSafedQuery[T], error) {
	q, err := c.untyped.QueryNested(parent, target)
	if err != nil {
		return nil, err
	}
	return &TypeSafedQuery[T]{
		untyped: q,
	}, nil
}

// Iter runs query and calls callback for each document
func (q *TypeSafedQuery[T]) Iter(ctx context.Context, f func(o *T) error) error {
	return q.untyped.Iter(ctx, func(o any) error {
		return f(o.(*T))
	})
}

// GetAll runs query and retrieve all results
// results are stored to `target` passed in `Query()`, `QueryGroup()` or `QueryNested()`
func (q *TypeSafedQuery[T]) GetAll(ctx context.Context) error {
	return q.untyped.GetAll(ctx)
}

// Where sets the condition for the query
func (q *TypeSafedQuery[T]) Where(path, op string, value interface{}) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.Where(path, op, value),
	}
}

// OrderBy sets the order of the query result
func (q *TypeSafedQuery[T]) OrderBy(path string, dir firestore.Direction) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.OrderBy(path, dir),
	}
}

// Offset sets the offset of the query
func (q *TypeSafedQuery[T]) Offset(n int) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.Offset(n),
	}
}

// Limit sets the max count of documents of the query
func (q *TypeSafedQuery[T]) Limit(n int) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.Limit(n),
	}
}

// LimitToLast sets the max count of documents of the query
func (q *TypeSafedQuery[T]) LimitToLast(n int) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.LimitToLast(n),
	}
}

// StartAt sets the start position of the query
func (q *TypeSafedQuery[T]) StartAt(docSnapshotOrFieldValues ...interface{}) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.StartAt(docSnapshotOrFieldValues...),
	}
}

// StartAt sets the start position of the query
func (q *TypeSafedQuery[T]) StartAfter(docSnapshotOrFieldValues ...interface{}) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.StartAfter(docSnapshotOrFieldValues...),
	}
}

// EndAt sets the end position of the query
func (q *TypeSafedQuery[T]) EndAt(docSnapshotOrFieldValues ...interface{}) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.EndAt(docSnapshotOrFieldValues...),
	}
}

// EndBefore set the end position of the query
func (q *TypeSafedQuery[T]) EndBefore(docSnapshotOrFieldValues ...interface{}) *TypeSafedQuery[T] {
	return &TypeSafedQuery[T]{
		untyped: q.untyped.EndBefore(docSnapshotOrFieldValues...),
	}
}
