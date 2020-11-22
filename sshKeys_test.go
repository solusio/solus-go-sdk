package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHKeysService_Create(t *testing.T) {
	data := SSHKeyCreateRequest{
		Name:   "name",
		Body:   "body",
		UserID: 1,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeSSHKey)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).SSHKeys.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeSSHKey, actual)
}

func TestSSHKeysService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ssh_keys/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).SSHKeys.Delete(context.Background(), 10)
	require.NoError(t, err)
}
