package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOsImageVersionsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_image_versions/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeOsImageVersion)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImageVersions.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeOsImageVersion, actual)
}

func TestOsImageVersionsService_Update(t *testing.T) {
	data := OsImageVersionRequest{
		Position:         1.5,
		Version:          "version",
		URL:              "http://foo/bar",
		CloudInitVersion: CloudInitVersionV0,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_image_versions/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeOsImageVersion)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).OsImageVersions.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeOsImageVersion, actual)
}

func TestOsImageVersionsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/os_image_versions/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).OsImageVersions.Delete(context.Background(), 10)
	require.NoError(t, err)
}
