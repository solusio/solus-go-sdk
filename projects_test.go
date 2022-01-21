package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectsService_Create(t *testing.T) {
	data := ProjectRequest{
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
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterProjects{}).ByName("name")

	actual, err := createTestClient(t, s.URL).Projects.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestProjectsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeProject)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeProject, actual)
}

func TestProjectService_Update(t *testing.T) {
	data := ProjectRequest{
		Name:        "name",
		Description: "description",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeProject)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeProject, actual)
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
