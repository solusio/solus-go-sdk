package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsService_Create(t *testing.T) {
	data := ProjectCreateRequest{
		Name:        "fake name",
		Description: "fake description",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeProject)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeProject, actual)
}

func TestProjectsService_List(t *testing.T) {
	expected := ProjectsResponse{
		Data: []Project{
			fakeProject,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestProjectsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Projects.Delete(context.Background(), 10)
	require.NoError(t, err)
}
