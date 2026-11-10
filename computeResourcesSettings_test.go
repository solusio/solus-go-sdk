package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResourcesService_SettingsGet(t *testing.T) {
	expected := ComputeResourceSettings{
		CachePath:     "fake cache path",
		ISOPath:       "fake ISO path",
		BackupTmpPath: "fake backup tmp path",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/settings", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.SettingsGet(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_SettingsUpdate(t *testing.T) {
	data := ComputeResourceSettings{
		CachePath:           "fake cache path",
		ISOPath:             "fake ISO path",
		BackupTmpPath:       "fake backup tmp path",
		VNCProxyPort:        1337,
		Limits:              ComputeResourceSettingsLimits{},
		Network:             ComputeResourceSettingsNetwork{},
		BalanceStrategy:     ComputeResourceBalanceStrategyMostStorageAvailable,
		VirtualizationTypes: []VirtualizationType{VirtualizationTypeKVM, VirtualizationTypeVZ},
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
