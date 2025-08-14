package simplestoretest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/ikedam/simplestore"
	"github.com/stretchr/testify/require"
)

// ClearFirestore clears the firestore database
// This works only when you use firestore emulator.
// This clears only the connecting database, not the whole emulator.
// See https://cloud.google.com/firestore/native/docs/emulator#clear_emulator_data
func ClearFirestore(t *testing.T, client *simplestore.Client) {
	addr := os.Getenv("FIRESTORE_EMULATOR_HOST")
	require.NotEmpty(t, addr, "FIRESTORE_EMULATOR_HOST is not set: ClearFirestore works only when you use firestore emulator")
	projectID := client.ProjectID
	if projectID == "" {
		projectID = "dummy-emulator-firestore-project"
		t.Logf("projectID is not set: assume `%s` for projectID", projectID)
	}
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"http://%s/emulator/v1/projects/%s/databases/%s/documents",
			addr,
			projectID,
			client.DatabaseID,
		),
		nil,
	)
	require.NoError(t, err, "failed to prepare request for ClearFirestore()")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "failed to clear firestore database in ClearFirestore()")
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "failed to clear firestore database in ClearFirestore()")
}
