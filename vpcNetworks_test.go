package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVPCNetworksService_Attach(t *testing.T) {
	req := VPCNetworkAttachRequest{
		ServerIDs: []int{1, 2},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/10/attach", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).VPCNetworks.Attach(context.Background(), 10, req)
	require.NoError(t, err)
}

func TestVPCNetworksService_Detach(t *testing.T) {
	req := VPCNetworkDetachRequest{
		ServerIDs: []int{3, 4},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/11/detach", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		w.WriteHeader(http.StatusOK)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).VPCNetworks.Detach(context.Background(), 11, req)
	require.NoError(t, err)
}
