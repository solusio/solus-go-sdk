package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServersService_List(t *testing.T) {
	expected := ServersResponse{
		Data: []Server{
			fakeServer,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[user_id]":             []string{"1"},
			"filter[compute_resource_id]": []string{"2"},
			"filter[status]":              []string{"status"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterServers{}).
		ByUserID(1).
		ByComputeResourceID(2).
		ByStatus("status")

	actual, err := createTestClient(t, s.URL).Servers.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestServersService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeServer, actual)
}

func TestServersService_Patch(t *testing.T) {
	data := ServerUpdateRequest{
		Name:        "name",
		BootMode:    BootModeRescue,
		Description: "description",
		UserData:    "data",
		FQDNs:       []string{"foo.example.com"},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/42", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Patch(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeServer, actual)
}

func TestServersService_Start(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/start", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Start(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestServersService_Stop(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/stop", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Stop(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestServersService_Restart(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/restart", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Restart(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestServersService_Backup(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/backups", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusCreated, fakeBackup)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Backup(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeBackup, actual)
}

func TestServersService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Servers.Delete(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}
