package simplestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type Query struct {
	q  firestore.Query
	tb *targetBuilder
}

// QuerySafe starts a new query for target
// target is also used as destination of `GetAll()`.
func (c *Client) QuerySafe(target *[]*any) (*Query, error) {
	collection, tb, err := c.getCollectionRef(target)
	if err != nil {
		return nil, err
	}
	return &Query{
		q:  collection.Query,
		tb: tb,
	}, nil
}

// Query starts a new query for target
// target is also used as destination of `GetAll()`.
// Panic if inappropriate target is specified.
func (c *Client) Query(target *[]*any) *Query {
	q, err := c.QuerySafe(target)
	if err != nil {
		panic(err)
	}
	return q
}

// QueryGroupSafe starts a new query for target as collection group
// target is also used as destination of `GetAll()`.
func (c *Client) QueryGroupSafe(target *[]*any) (*Query, error) {
	cgroup, tb, err := c.getCollectionGroupRef(target)
	if err != nil {
		return nil, err
	}
	return &Query{
		q:  cgroup.Query,
		tb: tb,
	}, nil
}

// QueryGroup starts a new query for target as collection group
// target is also used as destination of `GetAll()`.
// Panic if inappropriate target is specified.
func (c *Client) QueryGroup(target *[]*any) *Query {
	q, err := c.QueryGroupSafe(target)
	if err != nil {
		panic(err)
	}
	return q
}

// QueryNested starts a new query for target under the document specified by `parent`
// target is also used as destination of `GetAll()`.
func (c *Client) QueryNested(parent *any, target *[]*any) (*Query, error) {
	cgroup, tb, err := c.getNestedCollectionRef(parent, target)
	if err != nil {
		return nil, err
	}
	return &Query{
		q:  cgroup.Query,
		tb: tb,
	}, nil
}

// Iter runs query and calls callback for each document
func (q *Query) Iter(ctx context.Context, f func(o *any) error) error {
	iter := q.q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		dst := q.tb.createElement()
		err = doc.DataTo(dst)
		if err != nil {
			return err
		}
		err = f(dst.(*any))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAll runs query and retrieve all results
// results are stored to `target` passed in `Query()`, `QueryGroup()` or `QueryNested()`
func (q *Query) GetAll(ctx context.Context) error {
	return q.Iter(ctx, func(o *any) error {
		q.tb.append(o)
		return nil
	})
}

// Where sets the condition for the query
func (q *Query) Where(path, op string, value interface{}) *Query {
	newQ := *q
	newQ.q = q.q.Where(path, op, value)
	return &newQ
}

// OrderBy sets the order of the query result
func (q *Query) OrderBy(path string, dir firestore.Direction) *Query {
	newQ := *q
	newQ.q = q.q.OrderBy(path, dir)
	return &newQ
}

// Offset sets the offset of the query
func (q *Query) Offset(n int) *Query {
	newQ := *q
	newQ.q = q.q.Offset(n)
	return &newQ
}

// Limit sets the max count of documents of the query
func (q *Query) Limit(n int) *Query {
	newQ := *q
	newQ.q = q.q.Limit(n)
	return &newQ
}

// LimitToLast sets the max count of documents of the query
func (q *Query) LimitToLast(n int) *Query {
	newQ := *q
	newQ.q = q.q.LimitToLast(n)
	return &newQ
}

// StartAt sets the start position of the query
func (q *Query) StartAt(docSnapshotOrFieldValues ...interface{}) *Query {
	newQ := *q
	newQ.q = q.q.StartAt(docSnapshotOrFieldValues...)
	return &newQ
}

// StartAt sets the start position of the query
func (q *Query) StartAfter(docSnapshotOrFieldValues ...interface{}) *Query {
	newQ := *q
	newQ.q = q.q.StartAfter(docSnapshotOrFieldValues...)
	return &newQ
}

// EndAt sets the end position of the query
func (q *Query) EndAt(docSnapshotOrFieldValues ...interface{}) *Query {
	newQ := *q
	newQ.q = q.q.EndAt(docSnapshotOrFieldValues...)
	return q
}

// EndBefore set the end position of the query
func (q *Query) EndBefore(docSnapshotOrFieldValues ...interface{}) *Query {
	newQ := *q
	newQ.q = q.q.EndBefore(docSnapshotOrFieldValues...)
	return q
}
