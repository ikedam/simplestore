package simplestore

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestSimpleDoc struct {
	ID string
}

type TestChildDoc struct {
	Parent *TestSimpleDoc
	ID     string
}

type TestGrandChildDoc struct {
	Parent *TestChildDoc
	ID     string
}

type MyDocument struct {
	ID   string
	Name string
}

type ParentDocument struct {
	ID   string
	Name string
}

type ChildDocument struct {
	Parent *ParentDocument
	ID     string
	Name   string
}

func ptrOf[T any](v T) *T {
	return &v
}

type envMocker struct {
	origs map[string]*string
}

func newEnvMocker(envs map[string]*string) *envMocker {
	m := &envMocker{
		origs: make(map[string]*string, len(envs)),
	}
	for key, val := range envs {
		orig, ok := os.LookupEnv(key)
		if ok {
			m.origs[key] = &orig
		} else {
			m.origs[key] = nil
		}
		if val != nil {
			os.Setenv(key, *val)
		} else {
			os.Unsetenv(key)
		}
	}
	return m
}

func (m *envMocker) restore() {
	for key, val := range m.origs {
		if val != nil {
			os.Setenv(key, *val)
		} else {
			os.Unsetenv(key)
		}
	}
}

func mockEnvs(envs map[string]*string, f func()) {
	m := newEnvMocker(envs)
	f()
	m.restore()
}

type MockFirestore struct {
	m *envMocker
}

func (m *MockFirestore) SetupSuite() {
	// if FIRESTORE_EMULATOR_HOST is already set, do nothing
	if _, ok := os.LookupEnv("FIRESTORE_EMULATOR_HOST"); ok {
		return
	}
	envs := map[string]*string{
		"FIRESTORE_EMULATOR_HOST": ptrOf("localhost:8080"),
		"CLOUDSDK_CORE_PROJECT":   ptrOf("testproject"),
	}
	m.m = newEnvMocker(envs)
}

func (m *MockFirestore) TeardownSuite() {
	if m.m != nil {
		m.m.restore()
	}
}

func assertProgrammingError(tt assert.TestingT, err error, i ...interface{}) bool {
	return assert.ErrorAs(tt, err, ptrOf(&ProgrammingError{}))
}

func getDocumentOnlyPath(path string) string {
	pathList := strings.Split(path, "/")
	if len(pathList) <= 5 {
		return path
	}
	return strings.Join(pathList[5:], "/")
}

// clearAllDocuments deletes all documents in the specified collection.
// document should be specified with `&Document{}`.
func clearAllDocuments(t *testing.T, doc any) {
	ctx := context.Background()
	batchSize := 100

	accessor, err := newAccessor(reflect.TypeOf(doc))
	require.NoError(t, err)
	client, err := firestore.NewClient(context.Background(), getProjectID())
	require.NoError(t, err)
	err = DeleteCollection(ctx, client, accessor.collectionName, batchSize)
	require.NoError(t, err)
}

func DeleteCollection(ctx context.Context, client *firestore.Client, collectionPath string, batchSize int) error {
	collectionRef := client.Collection(collectionPath)
	for {
		// Get a batch of documents
		var documents []*firestore.DocumentSnapshot
		err := client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
			var err error
			documents, err = t.Documents(collectionRef.Limit(batchSize)).GetAll()
			return err
		})
		if err != nil {
			return err
		}

		// If no documents are left, break the loop
		if len(documents) == 0 {
			break
		}

		// Delete each document and its subcollections
		for _, doc := range documents {
			// Delete subcollections recursively
			subcollections, err := doc.Ref.Collections(ctx).GetAll()
			if err != nil {
				return err
			}
			for _, subcollection := range subcollections {
				err := DeleteCollection(
					ctx,
					client,
					fmt.Sprintf(
						"%v/%v/%v",
						collectionPath,
						doc.Ref.ID,
						subcollection.ID,
					),
					batchSize,
				)
				if err != nil {
					return err
				}
			}

			// Delete the document
			_, err = doc.Ref.Delete(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
