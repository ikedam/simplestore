package simplestoretest

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ikedam/simplestore"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// FirestoreTestSuite is a test suite that provides following features:
// - Simplestore client with a unique database ID
//   - Useful to avoid conflicts between test packages when using Firestore Emulator
//
// - Auto cleanup for each tests
type FirestoreTestSuite struct {
	suite.Suite
	SimplestoreClient   *simplestore.Client
	FirestoreDatabaseID string
}

// SetupSuite initializes the test suite with a unique database ID and simplestore client
func (s *FirestoreTestSuite) SetupSuite() {
	ctx := context.Background()

	// Generate a unique database ID to avoid conflicts between test packages
	s.FirestoreDatabaseID = generateUniqueDatabaseID(s.T())

	// Create a new client with the unique database ID
	client, err := simplestore.NewClientWithDatabase(
		ctx,
		s.FirestoreDatabaseID,
	)
	s.Require().NoError(err, "Failed to create simplestore client")

	s.SimplestoreClient = client
	ClearFirestore(s.T(), s.SimplestoreClient)
}

// TearDownSuite cleans up the test suite resources
func (s *FirestoreTestSuite) TearDownSuite() {
	if s.SimplestoreClient != nil {
		s.SimplestoreClient.Close()
	}
}

func (s *FirestoreTestSuite) TearDownTest() {
	ClearFirestore(s.T(), s.SimplestoreClient)
}

// generateUniqueDatabaseID generates a unique database ID for testing
// This prevents conflicts between different test packages running in parallel
func generateUniqueDatabaseID(t *testing.T) string {
	// Generate 8 random bytes and convert to hex string
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	require.NoError(t, err, "failed to generate random bytes")
	return fmt.Sprintf("test_%x", bytes)
}
