package solus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func TestVirtualServersService_StopWithRequest(t *testing.T) {
	req := VirtualServerStopRequest{
		Force: true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/stop", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.StopWithRequest(context.Background(), 10, req)
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

func TestVirtualServersService_Suspend(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/suspend", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Suspend(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_Resume(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/resume", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Resume(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_ResetUsage(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/reset_usage", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.ResetUsage(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_InstallGuestTools(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/install_guest_tools", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.InstallGuestTools(context.Background(), 10)
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
	data := VirtualServerResizeRequest{
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

func TestVirtualServersService_Reinstall(t *testing.T) {
	data := VirtualServerReinstallRequest{
		OSImageVersionID: 10,
		ApplicationID:    20,
		ApplicationData: map[string]interface{}{
			"domain": "example.com",
		},
		SSHKeys:  []int{1, 2},
		UserData: "some user data",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/reinstall", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeVirtualServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Reinstall(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeVirtualServer, actual)
}

func TestVirtualServersService_ResetPassword(t *testing.T) {
	sendToCurrentUser := true
	data := VirtualServerResetPasswordRequest{
		SendPasswordToCurrentUser: &sendToCurrentUser,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/reset_password", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeJSON(t, w, http.StatusOK, struct {
			Data struct {
				Password string `json:"password"`
			} `json:"data"`
		}{
			Data: struct {
				Password string `json:"password"`
			}{
				Password: "new-password",
			},
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.ResetPassword(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, "new-password", actual)
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

func TestVirtualServersService_Snapshots(t *testing.T) {
	expected := []Snapshot{fakeSnapshot}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/snapshots", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, struct {
			Data []Snapshot `json:"data"`
		}{
			Data: expected,
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Snapshots(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestVirtualServersService_Backups(t *testing.T) {
	expected := []Backup{fakeBackup}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/backups", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, struct {
			Data []Backup `json:"data"`
		}{
			Data: expected,
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Backups(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestVirtualServersService_Limits(t *testing.T) {
	expected := []VirtualServerLimit{
		{
			Name: "network_outgoing_traffic",
			Limit: struct {
				Limit     int    `json:"limit"`
				IsEnabled bool   `json:"is_enabled"`
				Unit      string `json:"unit"`
			}{
				Limit:     1000,
				IsEnabled: true,
				Unit:      "GiB",
			},
			Used:       123,
			IsExceeded: false,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/limits", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, struct {
			Data []VirtualServerLimit `json:"data"`
		}{
			Data: expected,
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.Limits(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestVirtualServersService_IPsCreate(t *testing.T) {
	delayed := false
	data := VirtualServerAdditionalIPsCreateRequest{
		Count:   2,
		Type:    IPv4,
		Delayed: &delayed,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.IPsCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_IPsUpdate(t *testing.T) {
	data := VirtualServerAdditionalIPsUpdateRequest{
		PrimaryID: 42,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.IPsUpdate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_IPsDelete(t *testing.T) {
	delayed := true
	data := VirtualServerAdditionalIPsDeleteRequest{
		IDs:     []int{1, 2},
		Delayed: &delayed,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assertRequestBody(t, r, data)
		writeJSON(t, w, http.StatusOK, struct {
			Data struct {
				IDs []int `json:"ids"`
			} `json:"data"`
		}{
			Data: struct {
				IDs []int `json:"ids"`
			}{
				IDs: []int{1, 2},
			},
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.IPsDelete(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2}, actual)
}

func TestVirtualServersService_IPDelete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/ips/15", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.IPDelete(context.Background(), 10, 15)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_DisksCreate(t *testing.T) {
	data := AdditionalDiskCreateRequest{
		Name:    "disk-1",
		Size:    20,
		OfferID: 7,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/disks", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusCreated, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.DisksCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_DiskUpdate(t *testing.T) {
	data := AdditionalDiskUpdateRequest{Name: "renamed-disk"}
	expected := []Disk{{
		ID:   1,
		Name: "renamed-disk",
	}}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/disks/2", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)
		writeJSON(t, w, http.StatusOK, struct {
			Data []Disk `json:"data"`
		}{
			Data: expected,
		})
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.DiskUpdate(context.Background(), 10, 2, data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestVirtualServersService_DiskDelete(t *testing.T) {
	delayed := false
	data := AdditionalDiskDeleteRequest{
		Delayed: &delayed,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/disks/2", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.DiskDelete(context.Background(), 10, 2, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServersService_DiskResize(t *testing.T) {
	delayed := false
	data := AdditionalDiskResizeRequest{
		Size:    100,
		Delayed: &delayed,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/servers/10/disks/2/resize", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)
		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VirtualServers.DiskResize(context.Background(), 10, 2, data)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}

func TestVirtualServerReadHardwareSettings(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/servers/42", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, `{"data":{"id":42,"compute_resource":{"id":15,"name":"compute","host":"compute.example.com","capabilities":{"kvm":true,"os_name":"linux"},"is_in_maintenance":true,"storages":[{"id":7}]},"settings":{"mac_address":"02:00:00:00:00:01","firmware":"efi","disk_cache_mode":"none","disk_driver":"virtio"}}}`)
		assert.NoError(t, err)
	})
	defer s.Close()
	vm, err := createTestClient(t, s.URL).VirtualServers.Get(context.Background(), 42)
	require.NoError(t, err)
	assert.Equal(t, ShortComputeResource{
		ID: 15, Name: "compute", Host: "compute.example.com",
		Capabilities:    map[string]interface{}{"kvm": true, "os_name": "linux"},
		IsInMaintenance: true, Storages: []ShortStorage{{ID: 7}},
	}, vm.ComputeResource)
	assert.Equal(t, VirtualServerSettings{
		MacAddress: "02:00:00:00:00:01", Firmware: "efi", DiskCacheMode: "none", DiskDriver: "virtio",
	}, vm.Settings)
}

func TestVirtualServerPatchFieldsHTTP(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusUnprocessableEntity} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPatch, r.Method)
				assert.Equal(t, "/servers/42", r.URL.Path)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"description":"","user_data":"#cloud-config","fqdns":[]}`, string(body))
				writeResponse(t, w, status, VirtualServer{ID: 42})
			})
			defer s.Close()
			empty, names, userData := "", []string{}, "#cloud-config"
			vm, err := createTestClient(t, s.URL).VirtualServers.PatchFields(context.Background(), 42,
				VirtualServerPatchRequest{Description: &empty, UserData: &userData, FQDNs: &names})
			if status == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, 42, vm.ID)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestVirtualServerUpdateSettingsTaskHTTP(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusUnprocessableEntity} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPatch, r.Method)
				assert.Equal(t, "/servers/42/settings", r.URL.Path)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"firmware":"efi"}`, string(body))
				writeResponse(t, w, status, Task{ID: 123, Status: TaskStatusQueued})
			})
			defer s.Close()
			task, err := createTestClient(t, s.URL).VirtualServers.UpdateSettingsTask(context.Background(), 42,
				VirtualServerUpdateSettingsRequest{Firmware: "efi"})
			if status == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, 123, task.ID)
				assert.False(t, task.IsFinished())
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestVirtualServerPatchBackupSettingsJSON(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(strconv.FormatBool(enabled), func(t *testing.T) {
			empty, names := "", []string{}
			raw, err := json.Marshal(VirtualServerPatchRequest{
				Description: &empty, FQDNs: &names,
				BackupSettings: &VirtualServerBackupSettings{
					Enabled: enabled,
					Schedule: VirtualServerBackupSettingsSchedule{
						Type: ServerBackupSettingsScheduleTypeWeekly,
						Time: VirtualServerBackupSettingsScheduleTime{Hour: 3, Minutes: 30},
						Days: []int{1, 3},
					},
					Limit: UnitPlanLimit{IsEnabled: true, Limit: 5, Unit: PlanLimitUnits},
				},
			})
			require.NoError(t, err)
			assert.JSONEq(t, fmt.Sprintf(`{
                "description":"", "fqdns":[],
                "backup_settings":{
                    "enabled":%t,
                    "schedule":{"type":"weekly","time":{"hour":3,"minutes":30},"days":[1,3]},
                    "limit":{"is_enabled":true,"limit":5,"unit":"units"}
                }
            }`, enabled), string(raw))
		})
	}
}

func TestVirtualServerPatchOmittedFieldsJSON(t *testing.T) {
	raw, err := json.Marshal(VirtualServerPatchRequest{})
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(raw))
}

func TestVirtualServersEmbeddedStorageCredentials(t *testing.T) {
	for _, credentials := range []string{"[]", "{}", "null"} {
		for _, operation := range []string{"Get", "List", "Patch", "PatchFields"} {
			t.Run(operation+"/"+credentials, func(t *testing.T) {
				vmJSON := fmt.Sprintf(`{"id":42,"compute_resource":{"id":15,"storages":[{"id":7,"name":null,"type":{"name":"lvm","formats":["raw"]},"status":"active","credentials":%s}]}}`, credentials)
				s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					body := `{"data":` + vmJSON + `}`
					if operation == "List" {
						assert.Equal(t, http.MethodGet, r.Method)
						assert.Equal(t, "/servers", r.URL.Path)
						body = `{"data":[{"id":1},` + vmJSON + `]}`
					} else {
						assert.Equal(t, "/servers/42", r.URL.Path)
						if operation == "Get" {
							assert.Equal(t, http.MethodGet, r.Method)
						} else {
							assert.Equal(t, http.MethodPatch, r.Method)
						}
					}
					_, err := io.WriteString(w, body)
					assert.NoError(t, err)
				})
				defer s.Close()
				service := createTestClient(t, s.URL).VirtualServers
				var vm VirtualServer
				var err error
				switch operation {
				case "Get":
					vm, err = service.Get(context.Background(), 42)
				case "List":
					var page VirtualServersResponse
					page, err = service.List(context.Background(), &FilterVirtualServers{})
					require.NoError(t, err)
					require.Len(t, page.Data, 2)
					vm = page.Data[1]
				case "Patch":
					vm, err = service.Patch(context.Background(), 42, VirtualServerUpdateRequest{Description: "updated"})
				case "PatchFields":
					description := "updated"
					vm, err = service.PatchFields(context.Background(), 42, VirtualServerPatchRequest{Description: &description})
				}
				require.NoError(t, err)
				assert.Equal(t, 42, vm.ID)
				require.Len(t, vm.ComputeResource.Storages, 1)
				storage := vm.ComputeResource.Storages[0]
				assert.Equal(t, 7, storage.ID)
				assert.Empty(t, storage.Name)
				assert.Equal(t, "active", storage.Status)
				assert.Equal(t, StorageType{Name: StorageTypeNameLVM, Formats: []ImageFormat{ImageFormatRaw}}, storage.Type)
			})
		}
	}
}
