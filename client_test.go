package simplestore

import (
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
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
