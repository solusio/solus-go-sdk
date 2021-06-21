package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResourcesService_Create(t *testing.T) {
	data := ComputerResourceCreateRequest{
		Name:      "name",
		Host:      "host",
		Login:     "login",
		Port:      1,
		Type:      "type",
		Password:  "password",
		Key:       "key",
		AgentPort: 2,
		IPBlocks:  []int{3, 4},
		Locations: []int{5, 6},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeComputeResource)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResource, actual)
}

func TestComputeResourcesService_Patch(t *testing.T) {
	data := ComputerResourceCreateRequest{
		Name:      "name",
		Host:      "host",
		Login:     "login",
		Port:      1,
		Type:      "type",
		Password:  "password",
		Key:       "key",
		AgentPort: 2,
		IPBlocks:  []int{3, 4},
		Locations: []int{5, 6},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/42", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeComputeResource)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Patch(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResource, actual)
}

func TestComputeResourcesService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeComputeResource)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResource, actual)
}

func TestComputeResourcesService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		assertRequestBody(t, r, deleteRequest{
			Force: true,
		})

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).ComputeResources.Delete(context.Background(), 10, true)
	require.NoError(t, err)
}

func TestComputeResourcesService_Networks(t *testing.T) {
	expected := []ComputeResourceNetwork{
		{
			ID:           "fake id",
			Name:         "fake network name",
			AddrConfType: "static",
			IPVersion:    4,
			IP:           "192.0.2.1",
			Mask:         "255.255.0.0",
			MaskSize:     16,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/networks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Networks(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_SetUpNetwork(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/setup_network", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, SetupNetworkRequest{
			ID:   "42",
			Type: ComputeResourceSettingsNetworkTypeBridged,
		})
	})
	defer s.Close()

	err := createTestClient(t, s.URL).ComputeResources.SetUpNetwork(context.Background(), 10, SetupNetworkRequest{
		ID:   "42",
		Type: ComputeResourceSettingsNetworkTypeBridged,
	})
	require.NoError(t, err)
}

func TestComputeResourcesService_PhysicalVolumes(t *testing.T) {
	expected := []ComputeResourcePhysicalVolume{
		{
			VGFree: "1",
			VGName: "2",
			VGSize: "3",
			PVUsed: "4",
		},
		{
			VGFree: "5",
			VGName: "6",
			VGSize: "7",
			PVUsed: "8",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/physical_volumes", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.PhysicalVolumes(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_ThinPools(t *testing.T) {
	expected := []ComputeResourceThinPool{
		{
			ConvertLV:       "1",
			CopyPercent:     "2",
			DataPercent:     "3",
			LVAttr:          "4",
			LVLayout:        "5",
			LVMetadataSize:  "6",
			LVName:          "7",
			LVSize:          "8",
			MetadataPrecent: "9",
			MirrorLog:       "10",
			MovePV:          "11",
			Origin:          "12",
			PoolLV:          "13",
			VGName:          "14",
		},
		{
			ConvertLV:       "15",
			CopyPercent:     "16",
			DataPercent:     "17",
			LVAttr:          "18",
			LVLayout:        "19",
			LVMetadataSize:  "20",
			LVName:          "21",
			LVSize:          "22",
			MetadataPrecent: "23",
			MirrorLog:       "24",
			MovePV:          "25",
			Origin:          "26",
			PoolLV:          "27",
			VGName:          "28",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/thin_pools", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.ThinPools(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
