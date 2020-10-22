package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestTasksService_List(t *testing.T) {
	expected := TasksResponse{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t,
			url.Values{
				"filter[action]":                 []string{"action"},
				"filter[status]":                 []string{"status"},
				"filter[compute_resource_id]":    []string{"1"},
				"filter[compute_resource_vm_id]": []string{"2"},
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

	f := (&FilterTasks{}).
		ByAction("action").
		ByStatus("status").
		ByComputeResourceID(1).
		ByComputeResourceVmID(2)

	p, err := c.Tasks.List(context.Background(), f)
	require.NoError(t, err)
	p.service = nil
	require.Equal(t, expected, p)
}
