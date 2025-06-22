# simplestore

Package simplestore provides a wrapper for `cloud.google.com/go/firestore`.
It handles collection names based on struct names.

## Performance

simpleclient utilizes reflection, the performance is not so good.

## Creating a Client

	ctx := context.Background()
	client, err := simplestore.New(ctx)
	if err != nil {
		// TODO: Handle error.
	}

The Google Cloud project name can be specified via environment variables `CLOUDSDK_CORE_PROJECT`, `GOOGLE_CLOUD_PROJECT`,
or determined from credentials.
You can specify the project name explicitly with `NewWithProjectID` .

## Struct for documents in simplestore

- simplestore treats a struct name as a collection name
- simplestore requires `ID` field in structs. `ID` field is treated as document id
- If a struct has a `Parent` field, it is treated as a parent document.
- Mapping between structs and documents is described in [the document of `cloud.google.com/go/firestore`]:https://pkg.go.dev/cloud.google.com/go/firestore#DocumentRef.Create

Example:

	type Document struct {
		ID    string
		Name  string
		Cache string `firestore:"-"`	// this field won't be stored in firestore.
	}

	// This document will be stored as /Document/123
	doc := &Docuemnt {
		ID:   "123",
		Name: "Alice",
	}

Example with parent:

	type ParentDocument struct {
		ID   string
		Name string
	}

	type ChildDocument struct {
		Parent *ParentDocument
		ID     string
		Name   string
	}

	// This document will be stored as /ParentDocument/c/ChildDocument/123
	doc := &ChildDocument {
		ParentDocument: &ParentDocument {
			ID:   "c",
			Name: "cryptography",
		},
		ID:   "123",
		Name: "Alice",
	}

	// This document will be stored as /ChildDocument/123
	doc := &ChildDocument {
		ParentDocument: nil,
		ID: "123",
		Name: "Alice",
	}

## GetDocumentID() / SetDocumentID()

simplestore treats `ID` field as document ID by default.
You can define `GetDocumentID()` and `SetDocumentID()` instead.

Interface:

	type IDer interface {
		GetDocumentID() string
		SetDocumentID(id string)
	}

Example Implementation:

	type CustomIDDocument struct {
		MyID string
		Name string
	}

	func (d *CustomIDDocument) GetDocumentID() string {
		if d.MyID == "" {
			return ""
		}
		// Suffix with fnv32 to avoid hotspot
		h := fnv.New32()
		h.Write([]byte(d.MyID))
		return fmt.Sprintf("%08x.%s", h.Sum32(), d.MyID)
	}

	func (d *CustomIDDocument) SetDocumentID(id string) {
		if id == "" {
			d.MyID = ""
			return
		}
		idx := strings.Index(id, ".")
		if s < 0 {
			// strange!
			d.MyID = ""
			return
		}
		d.MyID = id[s+1:]
	}


## Reading

For a simple document:

	// Retrieves document from /MyDocument/docid
	doc := &MyDocument {
		ID: "docid",
	}
	err := client.Get(ctx, doc)	// doc must be a pointer to a struct
	if err != nil {
		// TODO: Handle error.
	}
	fmt.Println(doc)

For multiple documents

	// Retrieves document from /MyDocument/docid1 and /MyDocument/docid2
	doc1 := &MyDocument {
		ID: "docid1",
	}
	doc2 := &MyDocument {
		ID: "docid2",
	}
	err := client.GetAll(ctx, []*MyDocument{doc1, doc2})
	if err != nil {
		// TODO: Handle error.
	}

## Writing

`Create` creates a new document, and returns an error if the document already exists:

	// ID is not set, and automatically generated.
	// You can specify ID to use manually instead.
	doc := &MyDocument {
		Name: "Alice",
	}
	_, err := client.Create(ctx, doc)
	if err != nil {
		// TODO: Handle error.
	}

`Set` creates a new document if not exist, or overwrite otherwise:

	doc := &MyDocument {
		ID: "123"
		Name: "Bob",
	}
	_, err := client.Update(ctx, doc)
	if err != nil {
		// TODO: Handle error.
	}

`Delete` deletes a document.

	_, err := client.Delete(ctx, doc)

## Queries

Start queries with `Query()`. Pass pointer to the slice:

	var docs []*MyDocument
	q := client.Query(*docs).OrderBy("name", firestore.Desc)

Following methods can be appllied:

* `Where()`
* `OrderBy()`
* `Offset()`
* `Limit()`
* `LimitToLast()`
* `StartAt()`
* `StartAfter()`
* `EndAt()`
* `EndBefore()`

To get documents, two ways are provided: `Iter()` and `GetAll()`.

With `Iter()`:

	err := q.Iter(ctx, func(o any) error {
		doc := o.(*MyDocument)
		fmt.Println(doc)
		// Returning error stops the iteration.
		// The error will be returned from Iter().
		return nil
	}
	if err != nil {
		// TODO: Handle error.
	}

With `GetAll()`:

	err := q.GetAll(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	for _, doc := range docs {
		fmt.Println(doc)
	}

For documents in subcollections:

	var parentDoc *ParentDocument = ...
	var childDocs []*ChildDocument
	q := client.QueryNested(parentDoc, &childDocs)

## Transactions

`RunTransaction` passes a new client for transaction.
You can call methods just like outside of transaction.

	err := client.RunTransaction(ctx, func(ctx context.Context, client *simplestore.Client) error {
		err := client.Get(ctx, doc)
		if err != nil {
			return err
		}
		doc.Name = "Bob"
		_, err = client.Set(ctx, doc)
		return err
	})
	if err != nil {
		// TODO: Handle error.
	}

## Type safed client

Many parameters of simpleclient.Client is typed `any`, and you can easily create runtime errors by passing unmached types.
You can use TypeSafedClient to avoid type assertion errors:

	doc := &MyDocument {
		ID: "docid",
	}
	err := TypeSafed[MyDocument](client).Get(ctx, doc)	// you can restrict to pass *MyDocument


# Table Mapping

simplestore supports table mapping to customize collection names and set readonly flags for specific structs.

## Basic Table Mapping

You can map struct names to custom collection names:

	client.AddTableMaps(map[string]string{
		"MyDocument": "custom_collection",
		"User":       "users",
	})

## Readonly Table Mapping

You can also set collections as readonly to prevent write operations:

	client.AddReadonlyTableMaps(map[string]string{
		"ReadOnlyDocument": "readonly_collection",
		"Config":           "configs",
	})

When a collection is marked as readonly, all write operations (Create, Set, Delete) will return an error.

## Example Usage

	type MyDocument struct {
		ID   string
		Name string
	}

	type ReadOnlyDocument struct {
		ID   string
		Name string
	}

	// Set up table mappings
	client.AddTableMaps(map[string]string{
		"MyDocument": "custom_collection",
	})
	client.AddReadonlyTableMaps(map[string]string{
		"ReadOnlyDocument": "readonly_collection",
	})

	// This will be stored in "custom_collection"
	doc := &MyDocument{Name: "Test"}
	_, err := client.Create(ctx, doc)

	// This will fail with readonly error
	readonlyDoc := &ReadOnlyDocument{Name: "Test"}
	_, err := client.Create(ctx, readonlyDoc) // Error: cannot create document in readonly collection
