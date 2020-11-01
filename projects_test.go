package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

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
