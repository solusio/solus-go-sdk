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

func TestComputeResourcesService_Create(t *testing.T) {
	expected := ComputeResource{}
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
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := ComputerResourceCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/compute_resources", r.URL.Path)
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

	l, err := c.ComputeResources.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}
