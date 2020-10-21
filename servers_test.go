package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestServersService_List(t *testing.T) {
	expected := ServersResponse{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/servers", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t,
			url.Values{
				"filter[user_id]":             []string{"1"},
				"filter[compute_resource_id]": []string{"2"},
				"filter[status]":              []string{"status"},
			}.Encode(),
			r.URL.Query().Encode(),
		)

		b, err := json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	c := createTestClient(t, s.URL)

	f := (&FilterServers{}).
		ByUserID(1).
		ByComputeResourceID(2).
		ByStatus("status")

	p, err := c.Servers.List(context.Background(), f)
	require.NoError(t, err)
	p.service = nil
	require.Equal(t, expected, p)
}
