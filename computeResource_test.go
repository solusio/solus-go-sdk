package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
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
		IpBlocks:  []int{3, 4},
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

func TestComputeResourcesService_Networks(t *testing.T) {
	expected := []ComputeResourceNetwork{
		{
			Id:           "fake id",
			Name:         "fake network name",
			AddrConfType: "static",
			IpVersion:    4,
			Ip:           "192.0.2.1",
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
		assertRequestBody(t, r, struct {
			Id string `json:"id"`
		}{
			Id: "42",
		})
	})
	defer s.Close()

	err := createTestClient(t, s.URL).ComputeResources.SetUpNetwork(context.Background(), 10, "42")
	require.NoError(t, err)
}
