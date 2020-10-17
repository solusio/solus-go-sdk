package solus

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
)

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
