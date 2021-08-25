package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPBlocksService_List(t *testing.T) {
	expected := IPBlocksResponse{
		Data: []IPBlock{
			fakeIPBlock,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterIPBlocks{}).ByName("name")

	actual, err := createTestClient(t, s.URL).IPBlocks.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestIPBlocksService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeIPBlock)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IPBlocks.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeIPBlock, actual)
}

func TestIPBlocksService_Update(t *testing.T) {
	data := IPBlockRequest{
		ComputeResources: []int{1, 2},
		Name:             "name",
		Type:             IPv4,
		Gateway:          "192.0.2.1",
		Ns1:              "192.0.2.2",
		Ns2:              "192.0.2.3",
		Netmask:          "255.255.255.0",
		From:             "192.0.2.10",
		To:               "192.0.2.240",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeIPBlock)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IPBlocks.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeIPBlock, actual)
}

func TestIPBlocksService_Create(t *testing.T) {
	data := IPBlockRequest{
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

		writeResponse(t, w, http.StatusCreated, fakeIPBlock)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IPBlocks.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeIPBlock, actual)
}

func TestIPBlocksService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).IPBlocks.Delete(context.Background(), 10)
	require.NoError(t, err)
}

func TestIPBlocksService_IPAddressCreate(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ip_blocks/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		writeResponse(t, w, http.StatusCreated, fakeIPBlockIPAddress)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).IPBlocks.IPAddressCreate(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeIPBlockIPAddress, actual)
}

func TestIPBlocksService_IPAddressDelete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ips/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).IPBlocks.IPAddressDelete(context.Background(), 10)
	require.NoError(t, err)
}
