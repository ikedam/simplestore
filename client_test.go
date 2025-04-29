package simplestore

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjectID(t *testing.T) {
	testcases := []struct {
		Name     string
		Envs     map[string]*string
		Expected string
	}{
		{
			Name: "None envs defined",
			Envs: map[string]*string{
				"CLOUDSDK_CORE_PROJECT": nil,
				"GOOGLE_CLOUD_PROJECT":  nil,
			},
			Expected: firestore.DetectProjectID,
		},
		{
			Name: "Both envs defined",
			Envs: map[string]*string{
				"CLOUDSDK_CORE_PROJECT": ptrOf("project1"),
				"GOOGLE_CLOUD_PROJECT":  ptrOf("project2"),
			},
			Expected: "project1",
		},
		{
			Name: "only later",
			Envs: map[string]*string{
				"CLOUDSDK_CORE_PROJECT": nil,
				"GOOGLE_CLOUD_PROJECT":  ptrOf("project2"),
			},
			Expected: "project2",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			mockEnvs(testcase.Envs, func() {
				actual := getProjectID()
				assert.Equal(t, testcase.Expected, actual)
			})
		})
	}
}
func TestNewClient(t *testing.T) {
	ctx := context.Background()

	client, err := New(ctx)
	require.NoError(t, err, "expected no error while creating client")
	assert.NotNil(t, client, "client should not be nil")
}

func TestNewWithProjectID(t *testing.T) {
	ctx := context.Background()
	projectID := "test-project"

	client, err := NewWithProjectID(ctx, projectID)
	require.NoError(t, err, "expected no error while creating client with project ID")
	assert.NotNil(t, client, "client should not be nil")
}
