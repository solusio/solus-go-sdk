package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestBackupNodesService_Create(t *testing.T) {
	data := BackupNodeRequest{
		Name:             "name",
		Type:             BackupNodeTypeSSHRsync,
		ComputeResources: []int{1, 2},
		Credentials: map[string]interface{}{
			"foo": "bar",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backup_nodes", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeBackupNode)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).BackupNodes.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeBackupNode, actual)
}

func TestBackupNodesService_Update(t *testing.T) {
	data := BackupNodeRequest{
		Name:             "name",
		Type:             BackupNodeTypeSSHRsync,
		ComputeResources: []int{1, 2},
		Credentials: map[string]interface{}{
			"foo": "bar",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backup_nodes/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeBackupNode)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).BackupNodes.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeBackupNode, actual)
}

func TestBackupNodesService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/backup_nodes/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).BackupNodes.Delete(context.Background(), 10)
	require.NoError(t, err)
}
