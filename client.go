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

package simplestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// KnownProjectIDEnvs is environment variables to specify project ID
var KnownProjectIDEnvs = []string{
	"CLOUDSDK_CORE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
}

// Client is a client for simplestore
// This wraps firestore client. You can get raw firestore client via `FirestoreClient`.
type Client struct {
	FirestoreClient             *firestore.Client
	FirestoreTransaction        *firestore.Transaction
	transactionFailureCallbacks []func()
}

// New returns a new client
// The project id can be configured via environment variables `CLOUDSDK_CORE_PROJECT`, `GOOGLE_CLOUD_PROJECT`
// or determined from credentials.
func New(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	return NewWithProjectID(ctx, getProjectID(), opts...)
}

// NewWithProjectID returns a new client
func NewWithProjectID(ctx context.Context, projectID string, opts ...option.ClientOption) (*Client, error) {
	client, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		FirestoreClient: client,
	}, nil
}

// Close cleans resource of this client
func (c *Client) Close() error {
	return c.FirestoreClient.Close()
}

// NewWithScope calls callback with new created client
// The client will be automatically closed.
func NewWithScope(ctx context.Context, f func(client *Client) error) error {
	client, err := New(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	return f(client)
}

func getProjectID() string {
	for _, env := range KnownProjectIDEnvs {
		value, ok := os.LookupEnv(env)
		if ok {
			return value
		}
	}
	return firestore.DetectProjectID
}
