package solus

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVirtualServerReadHardwareSettings(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/servers/42", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := io.WriteString(w, `{"data":{"id":42,"compute_resource":{"id":15},"settings":{"mac_address":"02:00:00:00:00:01","firmware":"efi","disk_cache_mode":"none","disk_driver":"virtio"}}}`)
		assert.NoError(t, err)
	})
	defer s.Close()
	vm, err := createTestClient(t, s.URL).VirtualServers.Get(context.Background(), 42)
	require.NoError(t, err)
	assert.Equal(t, 15, vm.ComputeResource.ID)
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
				assert.JSONEq(t, `{"description":"","user_data":"","fqdns":[]}`, string(body))
				writeResponse(t, w, status, VirtualServer{ID: 42})
			})
			defer s.Close()
			empty, names := "", []string{}
			vm, err := createTestClient(t, s.URL).VirtualServers.PatchFields(context.Background(), 42,
				VirtualServerPatchRequest{Description: &empty, UserData: &empty, FQDNs: &names})
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

func TestPlanUpdateAvailabilityJSON(t *testing.T) {
	for _, tt := range []struct {
		name    string
		ids     []int
		present bool
	}{
		{name: "omitted"}, {name: "cleared", ids: []int{}, present: true}, {name: "restricted", ids: []int{7}, present: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := json.Marshal(PlanUpdateRequest{AvailableLocations: tt.ids, AvailableOsImageVersions: tt.ids, AvailableApplications: tt.ids})
			require.NoError(t, err)
			var data map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(raw, &data))
			for _, key := range []string{"available_locations", "available_os_image_versions", "available_applications"} {
				value, present := data[key]
				assert.Equal(t, tt.present, present)
				if present {
					expected, err := json.Marshal(tt.ids)
					require.NoError(t, err)
					assert.JSONEq(t, string(expected), string(value))
				}
			}
		})
	}
}

func TestVirtualServerPatchEmptyAndDisabledJSON(t *testing.T) {
	empty := ""
	emptyList := []string{}
	raw, err := json.Marshal(VirtualServerPatchRequest{Description: &empty, UserData: &empty, FQDNs: &emptyList, BackupSettings: &VirtualServerBackupSettings{Enabled: false}})
	require.NoError(t, err)
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &data))
	assert.Equal(t, "", data["description"])
	assert.Equal(t, "", data["user_data"])
	assert.Equal(t, []interface{}{}, data["fqdns"])
	assert.Equal(t, false, data["backup_settings"].(map[string]interface{})["enabled"])
	assert.NotContains(t, data, "name")
}

func TestVirtualServerLegacyBackupPatchDoesNotDisableBackups(t *testing.T) {
	raw, err := json.Marshal(VirtualServerUpdateRequest{
		BackupSettings: &VirtualServerBackupSettings{Limit: UnitPlanLimit{Limit: 5}},
	})
	require.NoError(t, err)
	var data map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &data))
	assert.NotContains(t, data["backup_settings"].(map[string]interface{}), "enabled")
}
