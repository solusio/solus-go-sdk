package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackup_IsFinished(t *testing.T) {
	testCases := map[BackupStatus]bool{
		BackupStatusPending:    false,
		BackupStatusInProgress: false,
		BackupStatusFailed:     true,
		BackupStatusCreated:    true,
	}

	for status, expected := range testCases {
		status := status
		expected := expected
		t.Run(string(status), func(t *testing.T) {
			require.Equal(t, expected, Backup{Status: status}.IsFinished())
		})
	}
}

func TestBackupsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backups/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeBackup)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Backups.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeBackup, actual)
}

func TestBackupsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backups/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Backups.Delete(context.Background(), 10)
	require.NoError(t, err)
}

func TestBackupsService_Restore(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/backups/10/restore", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			writeResponse(t, w, http.StatusOK, fakeTask)
		})
		defer s.Close()

		actual, err := createTestClient(t, s.URL).Backups.Restore(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, fakeTask, actual)
	})

	t.Run("negative", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/backups/10/restore", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			writeResponse(t, w, http.StatusNotFound, fakeTask)
		})
		defer s.Close()

		_, err := createTestClient(t, s.URL).Backups.Restore(context.Background(), 10)
		assert.EqualError(t, err, "HTTP POST backups/10/restore returns 404 status code")
	})
}
