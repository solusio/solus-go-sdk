package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestIpBlocksService_List(t *testing.T) {
	expected := IpBlocksResponse{
		Data: []IpBlock{
			fakeIpBlock,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IpBlocks.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestIpBlocksService_Create(t *testing.T) {
	data := IpBlockCreateRequest{
		ComputeResources: []int{1, 2},
		Name:             "name",
		Type:             IPv4,
		Gateway:          "gateway",
		Ns1:              "ns1",
		Ns2:              "ns2",
		Netmask:          "netmask",
		From:             "from",
		To:               "to",
		Range:            "range",
		Subnet:           3,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeIpBlock)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IpBlocks.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeIpBlock, actual)
}

func TestIpBlocksService_IpAddressCreate(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusCreated, fakeIpBlockIpAddress)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IpBlocks.IpAddressCreate(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeIpBlockIpAddress, actual)
}

func TestIpBlocksService_IpAddressDelete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ips/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).IpBlocks.IpAddressDelete(context.Background(), 10)
	require.NoError(t, err)
}
