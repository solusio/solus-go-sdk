// Autogenerated file. Do not edit!

package solus

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectServersResponse_Next(t *testing.T) {
	page := int32(1)

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		p := atomic.LoadInt32(&page)

		assert.Equal(t, "/projectservers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, strconv.Itoa(int(p)), r.URL.Query().Get("page"))

		if p == 3 {
			writeJSON(t, w, http.StatusOK, ProjectServersResponse{Data: []Server{{ID: int(p)}}})
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
			Data: []Server{{ID: int(p)}},
		})
	})
	defer s.Close()

	resp := ProjectServersResponse{
		paginatedResponse: paginatedResponse{
			Links: ResponseLinks{
				Next: fmt.Sprintf("%s/projectservers?page=1", s.URL),
			},
			service: &service{createTestClient(t, s.URL)},
		},
	}

	i := 1
	for resp.Next(context.Background()) {
		require.Equal(t, []Server{{ID: i}}, resp.Data)
		i++
	}
	require.NoError(t, resp.err)
	require.Equal(t, 4, i, "Expects to get 3 entity, but got less")
}
