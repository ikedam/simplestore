package simplestore

import (
	"os"
	"strings"

	"github.com/stretchr/testify/assert"
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
