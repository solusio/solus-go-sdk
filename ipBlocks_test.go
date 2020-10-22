package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestIpBlocksService_Create(t *testing.T) {
	expected := IpBlock{}
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
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := IpBlockCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/ip_blocks", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, data, d)

		b, err = json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(201)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.IpBlocks.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}
