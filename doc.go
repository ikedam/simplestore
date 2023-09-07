/*
Copyright 2023 IKEDA Yasuyuki

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

/*
Package simplestore provides a wrapper for `cloud.google.com/go/firestore`.
It handles collection names based on struct names.

# Creating a Client

	ctx := context.Background()
	client, err := simplestore.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

The Google Cloud project name can be specified via environment variables `CLOUDSDK_CORE_PROJECT`, `GOOGLE_CLOUD_PROJECT`,
or determined from credentials.
You can specify the project name explicitly with `NewWithProjectID` .

# Struct for documents in simplestore

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

# Reading

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

# Writing

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
*/
package simplestore
