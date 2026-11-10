package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService_Notifications(t *testing.T) {
	expected := []Notification{
		{
			ID:   1,
			Type: "info",
			Data: map[string]interface{}{
				"message": "hello",
			},
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/notifications", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.Notifications(context.Background())
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestAccountService_NotificationsDelete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/notifications", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Account.NotificationsDelete(context.Background())
	require.NoError(t, err)
}

func TestAccountService_SSHKeys(t *testing.T) {
	expected := SSHKeysResponse{
		Data: []SSHKey{fakeSSHKey},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.SSHKeys(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestAccountService_SSHKeyCreate(t *testing.T) {
	data := AccountSSHKeyCreateRequest{
		Name: "test key",
		Body: "ssh-rsa AAA",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/ssh_keys", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, fakeSSHKey)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.SSHKeyCreate(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeSSHKey, actual)
}

func TestAccountService_Tokens(t *testing.T) {
	expected := APITokensResponse{
		Data: []APIToken{
			{
				ID:   "1",
				Name: "token",
			},
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/tokens", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.Tokens(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestAccountService_TokenCreate(t *testing.T) {
	data := AccountTokenCreateRequest{
		Name: "example token",
	}
	expected := AccessToken{
		Token: APIToken{
			ID:   "1",
			Name: "example token",
		},
		AccessToken: "token",
		TokenType:   "Bearer",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account/tokens", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.TokenCreate(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
