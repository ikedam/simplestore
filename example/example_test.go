package example

import (
	"context"
	"testing"

	"github.com/ikedam/simplestore/simplestoretest"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExampleTestSuite struct {
	simplestoretest.FirestoreTestSuite
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

type Document struct {
	ID   string
	Name string
}

func (s *ExampleTestSuite) TestSomething() {
	ctx := context.Background()
	doc := &Document{
		ID:   "123",
		Name: "Alice",
	}
	_, err := s.SimplestoreClient.Create(ctx, doc)
	s.Require().NoError(err)
	actualDoc := &Document{
		ID: "123",
	}
	err = s.SimplestoreClient.Get(ctx, actualDoc)
	s.Require().NoError(err)
	s.Assert().Equal(doc, actualDoc)
}

func (s *ExampleTestSuite) TestSomethingEnsureCleared() {
	ctx := context.Background()
	actualDoc := &Document{
		ID: "123",
	}
	// This results NotFound even after TestSomething()
	// as FirestoreTestSuite clears the database after each test.
	err := s.SimplestoreClient.Get(ctx, actualDoc)
	s.Require().Error(err)
	s.Assert().Equal(codes.NotFound, status.Code(err))
}

func (s *ExampleTestSuite) TestClearData() {
	ctx := context.Background()
	doc := &Document{
		ID:   "123",
		Name: "Alice",
	}
	_, err := s.SimplestoreClient.Create(ctx, doc)
	s.Require().NoError(err)
	actualDoc := &Document{
		ID: "123",
	}
	err = s.SimplestoreClient.Get(ctx, actualDoc)
	s.Require().NoError(err)
	s.Assert().Equal(doc, actualDoc)

	// clears the firestore database
	simplestoretest.ClearFirestore(s.T(), s.SimplestoreClient)

	err = s.SimplestoreClient.Get(ctx, actualDoc)
	s.Require().Error(err)
	s.Assert().Equal(codes.NotFound, status.Code(err))
}
