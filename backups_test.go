package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestBackup_IsFinished(t *testing.T) {
	testCases := map[BackupStatus]bool{
		BackupStatusPending:    false,
		BackupStatusInProgress: false,
		BackupStatusFailed:     true,
		BackupStatusCreated:    true,
	}

	for status, expected := range testCases {
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

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Backups.Delete(context.Background(), 10)
	require.NoError(t, err)
}

func TestBackupsService_Restore(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backups/10/restore", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Backups.Restore(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}
