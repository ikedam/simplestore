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

	type Document struct {
		Parent *ParentDocument
		ID     string
		Name   string
	}

	// This document will be stored as /ParentDocument/c/Document/123
	doc := &Docuemnt {
		ParentDocument: &ParentDocument {
			ID:   "c",
			Name: "cryptography",
		},
		ID:   "123",
		Name: "Alice",
	}

	// This document will be stored as /Document/123
	doc := &Docuemnt {
		ParentDocument: nil,
		ID: "123",
		Name: "Alice",
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
	found_docs, err := client.GetAll(ctx, []*MyDocument{doc1, doc2})
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

You can use SQL to select documents from a collection. Begin with the collection, and
build up a query using Select, Where and other methods of Query.

	q := states.Where("pop", ">", 10).OrderBy("pop", firestore.Desc)

Supported operators include '<', '<=', '>', '>=', '==', 'in', 'array-contains', and
'array-contains-any'.

Call the Query's Documents method to get an iterator, and use it like
the other Google Cloud Client iterators.

	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Println(doc.Data())
	}

To get all the documents in a collection, you can use the collection itself
as a query.

	iter = client.Collection("States").Documents(ctx)

## Collection Group Partition Queries

You can partition the documents of a Collection Group allowing for smaller subqueries.

	collectionGroup = client.CollectionGroup("States")
	partitions, err = collectionGroup.GetPartitionedQueries(ctx, 20)

You can also Serialize/Deserialize queries making it possible to run/stream the
queries elsewhere; another process or machine for instance.

	queryProtos := make([][]byte, 0)
	for _, query := range partitions {
		protoBytes, err := query.Serialize()
		// handle err
		queryProtos = append(queryProtos, protoBytes)
		...
	}

	for _, protoBytes := range queryProtos {
		query, err := client.CollectionGroup("").Deserialize(protoBytes)
		...
	}

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

# Type safed client

Many parameters of simpleclient.Client is typed `any`, and you can easily create runtime errors by passing unmached types.
You can use TypeSafedClient to avoid type assertion errors:

	doc := &MyDocument {
		ID: "docid",
	}
	err := TypeSafed[MyDocument](client).Get(ctx, doc)	// you can restrict to pass *MyDocument
