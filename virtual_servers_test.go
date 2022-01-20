package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVirtualServerService_Create(t *testing.T) {
	data := VirtualServerCreateRequest{
		Name:     "foo",
		BootMode: BootModeDisk,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeVirtualServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeVirtualServer, actual)
}

func TestVirtualServersService_setDefaultsForCreateRequest(t *testing.T) {
	cc := map[string]struct {
		given    VirtualServerCreateRequest
		expected VirtualServerCreateRequest
	}{
		"empty": {
			expected: VirtualServerCreateRequest{
				BootMode: BootModeDisk,
			},
		},

		"not empty": {
			given: VirtualServerCreateRequest{
				Name:             "name",
				BootMode:         BootModeRescue,
				Description:      "description",
				UserData:         "user data",
				FQDNs:            []string{"fqdns"},
				Password:         "password",
				SSHKeys:          []int{1},
				PlanID:           2,
				ProjectID:        3,
				LocationID:       4,
				OSImageVersionID: 5,
				ApplicationID:    6,
				ApplicationData:  map[string]interface{}{"foo": "bar"},
			},
			expected: VirtualServerCreateRequest{
				Name:             "name",
				BootMode:         BootModeRescue,
				Description:      "description",
				UserData:         "user data",
				FQDNs:            []string{"fqdns"},
				Password:         "password",
				SSHKeys:          []int{1},
				PlanID:           2,
				ProjectID:        3,
				LocationID:       4,
				OSImageVersionID: 5,
				ApplicationID:    6,
				ApplicationData:  map[string]interface{}{"foo": "bar"},
			},
		},
	}

	for name, c := range cc {
		t.Run(name, func(t *testing.T) {
			(&VirtualServersService{}).setDefaultsForCreateRequest(&c.given)
			assert.Equal(t, c.expected, c.given)
		})
	}
}

func TestVirtualServersService_List(t *testing.T) {
	expected := VirtualServersResponse{
		Data: []VirtualServer{
			fakeVirtualServer,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[user_id]":             []string{"1"},
			"filter[compute_resource_id]": []string{"2"},
			"filter[status]":              []string{"status"},
			"filter[virtualization_type]": []string{string(VirtualizationTypeKVM)},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterVirtualServers{}).
		ByUserID(1).
		ByComputeResourceID(2).
		ByStatus("status").
		ByVirtualizationType(VirtualizationTypeKVM)

	actual, err := createTestClient(t, s.URL).VirtualServers.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestVirtualServersService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeVirtualServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeVirtualServer, actual)
}

func TestVirtualServersService_Patch(t *testing.T) {
	data := VirtualServerUpdateRequest{
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

		writeResponse(t, w, http.StatusOK, fakeVirtualServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Patch(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeVirtualServer, actual)
}

func TestVirtualServersService_Start(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/start", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Start(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_Stop(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/stop", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Stop(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_Restart(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/restart", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Restart(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_Backup(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/servers/10/backups", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			writeResponse(t, w, http.StatusCreated, fakeBackup)
		})
		defer s.Close()

		actual, err := createTestClient(t, s.URL).VirtualServers.Backup(context.Background(), 10)
		require.NoError(t, err)
		require.Equal(t, fakeBackup, actual)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			asserter, addr := startBrokenTestServer(t)

			_, err := createTestClient(t, addr).VirtualServers.Backup(context.Background(), 10)
			asserter(t, http.MethodPost, "/servers/10/backups", err)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/servers/10/backups", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).VirtualServers.Backup(context.Background(), 10)
			assert.EqualError(t, err, "HTTP POST servers/10/backups returns 400 status code")
		})
	})
}

func TestVirtualServersService_resize(t *testing.T) {
	data := ViretualServerResizeRequest{
		PreserveDisk: true,
		PlanID:       42,
		BackupSettings: &VirtualServerBackupSettings{
			Enabled: true,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/resize", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Resize(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Delete(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_SnapshotsCreate(t *testing.T) {
	data := SnapshotRequest{
		Name: "name",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/snapshots", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeSnapshot)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.SnapshotsCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeSnapshot, actual)
}
