package solus

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSVersionHardwareSettingsRoundTrip(t *testing.T) {
	for _, data := range []string{
		`{"id":40,"vs_disk_driver":"virtio","vs_firmware":"efi"}`,
		`{"id":40,"vs_disk_driver":null,"vs_firmware":"bios"}`,
	} {
		var current OsImageVersion
		require.NoError(t, json.Unmarshal([]byte(data), &current))
		request := OsImageVersionRequest{VSDiskDriver: current.VSDiskDriver, VSFirmware: current.VSFirmware, IncludeVSDiskDriver: true}
		body, err := json.Marshal(request)
		require.NoError(t, err)
		var got map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &got))
		require.NotNil(t, current.VSFirmware)
		require.Equal(t, *current.VSFirmware, got["vs_firmware"])
		if current.VSDiskDriver != nil {
			require.Equal(t, *current.VSDiskDriver, got["vs_disk_driver"])
		} else {
			require.Contains(t, got, "vs_disk_driver")
			require.Nil(t, got["vs_disk_driver"])
		}
	}
}

func TestOSVersionHardwareSettingsOmittedByLegacyCallers(t *testing.T) {
	body, err := json.Marshal(OsImageVersionRequest{})
	require.NoError(t, err)
	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &got))
	require.NotContains(t, got, "vs_firmware")
	require.NotContains(t, got, "vs_disk_driver")
}
