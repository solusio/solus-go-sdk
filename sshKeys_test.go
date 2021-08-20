package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHKeysService_List(t *testing.T) {
	expected := SSHKeysResponse{
		Data: []SSHKey{
			fakeSSHKey,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterSSHKeys{}).ByName("name")

	actual, err := createTestClient(t, s.URL).SSHKeys.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestSSHKeysService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ssh_keys/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeSSHKey)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).SSHKeys.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeSSHKey, actual)
}

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
