package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/snapshots/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeSnapshot)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Snapshots.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeSnapshot, actual)
}

func TestSnapshotsService_Revert(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/snapshots/10/revert", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Snapshots.Revert(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestSnapshotsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/snapshots/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Snapshots.Delete(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}
