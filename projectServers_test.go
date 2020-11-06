package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestProjectsService_ServersCreate(t *testing.T) {
	data := ProjectServersCreateRequest{
		Name:             "name",
		PlanID:           1,
		LocationID:       2,
		OsImageVersionID: 3,
		SshKeys:          []int{4, 5},
		UserData:         "user data",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/42/servers", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.ServersCreate(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeServer, actual)
}

func TestProjectsService_ServersListAll(t *testing.T) {
	page := int32(0)

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		p := atomic.LoadInt32(&page)

		assert.Equal(t, "/projects/1/servers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		if page == 0 {
			assert.Equal(t, "", r.URL.Query().Get("page"))
		} else {
			assert.Equal(t, strconv.Itoa(int(p)), r.URL.Query().Get("page"))
		}

		if p == 2 {
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

	c := createTestClient(t, s.URL)

	actual, err := c.Projects.ServersListAll(context.Background(), 1)
	require.NoError(t, err)

	require.Equal(t, []Server{
		{ID: 0},
		{ID: 1},
		{ID: 2},
	}, actual)
}

func TestProjectsService_Servers(t *testing.T) {
	expected := ProjectServersResponse{
		Data: []Server{
			fakeServer,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/projects/42/servers", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Projects.Servers(context.Background(), 42)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}
