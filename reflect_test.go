package simplestore

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ReflectTestSuite struct {
	suite.Suite
	MockFirestore
}

func TestReflect(t *testing.T) {
	suite.Run(t, new(ReflectTestSuite))
}

func (s *ReflectTestSuite) TestNewAccessor() {
	t := s.T()
	testcases := []struct {
		Name  string
		Doc   any
		Error assert.ErrorAssertionFunc
	}{
		{
			Name: "Valid without parent",
			Doc:  &TestSimpleDoc{},
		},
		{
			Name: "Valid with parent",
			Doc:  &TestChildDoc{},
		},
		{
			Name: "Valid with grand parent",
			Doc:  &TestGrandChildDoc{},
		},
		{
			Name:  "Not a pointer",
			Doc:   TestSimpleDoc{},
			Error: assertProgrammingError,
		},
		{
			Name:  "A pointer of non-struct",
			Doc:   ptrOf("test"),
			Error: assertProgrammingError,
		},
		{
			Name:  "No ID Field",
			Doc:   &struct{}{},
			Error: assertProgrammingError,
		},
		{
			Name:  "ID Field is not a string",
			Doc:   &struct{ ID int64 }{},
			Error: assertProgrammingError,
		},
		{
			Name: "Parent is not a pointer",
			Doc: &struct {
				Parent TestSimpleDoc
				ID     string
			}{},
			Error: assertProgrammingError,
		},
		{
			Name: "Parent is bad for document",
			Doc: &struct {
				Parent *struct {
					ID int64
				}
				ID string
			}{},
			Error: assertProgrammingError,
		},
		{
			Name: "Grand parent is bad for document",
			Doc: &struct {
				Parent *struct {
					Parent *struct {
						ID int64
					}
					ID string
				}
				ID string
			}{},
			Error: assertProgrammingError,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			_, err := newAccessor(reflect.TypeOf(testcase.Doc))
			if testcase.Error != nil {
				testcase.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *ReflectTestSuite) TestGetDocumentRefWithoutParent() {
	t := s.T()
	testcases := []struct {
		Name         string
		Doc          *TestSimpleDoc
		MightNew     *bool
		ExpectedPath string
		ExpectNew    bool
		Error        assert.ErrorAssertionFunc
	}{
		{
			Name: "Existing doc",
			Doc: &TestSimpleDoc{
				ID: "123",
			},
			ExpectedPath: "TestSimpleDoc/123",
			ExpectNew:    false,
		},
	}
	ctx := context.Background()
	err := NewWithScope(ctx, func(c *Client) error {
		for _, testcase := range testcases {
			t.Run(testcase.Name, func(t *testing.T) {
				accessor, err := newAccessor(reflect.TypeOf(testcase.Doc))
				require.NoError(t, err)
				var mightNewList []bool
				if testcase.MightNew == nil {
					mightNewList = []bool{true, false}
				} else {
					mightNewList = []bool{*testcase.MightNew}
				}
				for _, mightNew := range mightNewList {
					doc, isNew, err := accessor.getDocumentRef(c, reflect.ValueOf(testcase.Doc), mightNew)
					if testcase.Error != nil {
						testcase.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
					assert.Equal(t, testcase.ExpectedPath, getDocumentOnlyPath(doc.Path))
					assert.Equal(t, testcase.ExpectNew, isNew)
				}
			})
		}
		return nil
	})
	require.NoError(t, err)
}
