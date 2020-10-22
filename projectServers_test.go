package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestProjectsService_ServersCreate(t *testing.T) {
	expected := Server{}
	data := ProjectServersCreateRequest{
		Name:             "name",
		PlanId:           1,
		LocationId:       2,
		OsImageVersionId: 3,
		SshKeys:          []int{4, 5},
		UserData:         "user data",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := ProjectServersCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/projects/42/servers", r.URL.Path)
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

	l, err := c.Projects.ServersCreate(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}

func TestProjectsService_ServersListAll(t *testing.T) {
	page := int32(0)

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		p := atomic.LoadInt32(&page)

		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/projects/1/servers", r.URL.Path)
		if page == 0 {
			require.Equal(t, "", r.URL.Query().Get("page"))
		} else {
			require.Equal(t, strconv.Itoa(int(p)), r.URL.Query().Get("page"))
		}

		if p == 2 {
			writeJSON(t, w, http.StatusOK, ProjectServersResponse{Data: []Server{{Id: int(p)}}})
			return
		}
		atomic.AddInt32(&page, 1)

		q := r.URL.Query()
		q.Set("page", strconv.Itoa(int(p)+1))
		r.URL.RawQuery = q.Encode()

		writeJSON(t, w, http.StatusOK, ProjectServersResponse{
			paginatedResponse: paginatedResponse{
				Links: ResponseLinks{
					Next: r.URL.String(),
				},
			},
			Data: []Server{{Id: int(p)}},
		})
	})
	defer s.Close()

	c := createTestClient(t, s.URL)

	servers, err := c.Projects.ServersListAll(context.Background(), 1)
	require.NoError(t, err)

	require.Equal(t, []Server{
		{Id: 0},
		{Id: 1},
		{Id: 2},
	}, servers)
}
