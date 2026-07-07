package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

var fakeVPCNetwork = VPCNetwork{
	ID:       1,
	Name:     "private network",
	ListType: VPCNetworkListTypeRange,
	From:     "10.0.0.10",
	To:       "10.0.0.20",
	Netmask:  "255.255.255.0",
	Location: ShortLocation{ID: fakeLocation.ID, Name: fakeLocation.Name},
	User:     fakeUser,
}

var fakeVPCNetworkIP = VPCNetworkIP{
	ID:         2,
	IP:         "10.0.0.11",
	Server:     fakeVirtualServer,
	VPCNetwork: ShortVPCNetwork{ID: fakeVPCNetwork.ID, Name: fakeVPCNetwork.Name},
}

func TestVPCNetworksService_Create(t *testing.T) {
	req := VPCNetworkCreateRequest{
		Name:       "private network",
		ListType:   VPCNetworkListTypeRange,
		From:       "10.0.0.10",
		To:         "10.0.0.20",
		Netmask:    "255.255.255.0",
		LocationID: 4,
		UserID:     null.IntFrom(3),
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusCreated, fakeVPCNetwork)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VPCNetworks.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, fakeVPCNetwork, actual)
}

func TestVPCNetworksService_List(t *testing.T) {
	expected := VPCNetworksResponse{
		Data: []VPCNetwork{fakeVPCNetwork},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"private"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterVPCNetworks{}).ByName("private")

	actual, err := createTestClient(t, s.URL).VPCNetworks.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestVPCNetworksService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeVPCNetwork)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VPCNetworks.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeVPCNetwork, actual)
}

func TestVPCNetworksService_Update(t *testing.T) {
	req := VPCNetworkUpdateRequest{
		Name:     "updated network",
		ListType: VPCNetworkListTypeSet,
		From:     "10.0.0.30",
		To:       "10.0.0.40",
		Netmask:  "255.255.255.0",
		UserID:   null.IntFrom(7),
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/12", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusOK, fakeVPCNetwork)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VPCNetworks.Update(context.Background(), 12, req)
	require.NoError(t, err)
	require.Equal(t, fakeVPCNetwork, actual)
}

func TestVPCNetworksService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/13", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).VPCNetworks.Delete(context.Background(), 13)
	require.NoError(t, err)
}

func TestVPCNetworksService_AddIP(t *testing.T) {
	req := VPCNetworkIPAddRequest{IP: "10.0.0.15"}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/14/add_ip", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusCreated, fakeVPCNetworkIP)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VPCNetworks.AddIP(context.Background(), 14, req)
	require.NoError(t, err)
	require.Equal(t, fakeVPCNetworkIP, actual)
}

func TestVPCNetworksService_AddIPs(t *testing.T) {
	req := VPCNetworkIPsAddRequest{
		IPs: []string{"10.0.0.15", "10.0.0.16"},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/15/add_ips", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)
		w.WriteHeader(http.StatusCreated)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).VPCNetworks.AddIPs(context.Background(), 15, req)
	require.NoError(t, err)
}

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

func TestVPCNetworksService_ListIPs(t *testing.T) {
	expected := VPCNetworkIPsResponse{
		Data: []VPCNetworkIP{fakeVPCNetworkIP},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/10/ips", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).VPCNetworks.ListIPs(context.Background(), 10)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestVPCNetworksService_DeleteIP(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/vpc_networks/10/ips/9", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).VPCNetworks.DeleteIP(context.Background(), 10, 9)
	require.NoError(t, err)
}
