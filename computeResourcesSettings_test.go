package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestComputeResourcesService_SettingsUpdate(t *testing.T) {
	data := ComputeResourceSettings{
		CachePath:       "fake cache path",
		ISOPath:         "fake ISO path",
		BackupTmpPath:   "fake backup tmp path",
		VNCProxyPort:    1337,
		Limits:          ComputeResourceSettingsLimits{},
		Network:         ComputeResourceSettingsNetwork{},
		BalanceStrategy: ComputeResourceBalanceStrategyMostSpaceAvailable,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/settings", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, data)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.SettingsUpdate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, data, actual)
}
