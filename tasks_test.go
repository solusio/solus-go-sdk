package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestTask_IsFinished(t *testing.T) {
	testCases := map[TaskStatus]bool{
		TaskStatusPending:        false,
		TaskStatusQueued:         false,
		TaskStatusRunning:        false,
		TaskStatusDone:           true,
		TaskStatusDoneWithErrors: true,
		TaskStatusFailed:         true,
		TaskStatusCanceled:       true,
	}

	for status, expected := range testCases {
		t.Run(string(status), func(t *testing.T) {
			require.Equal(t, expected, Task{Status: status}.IsFinished())
		})
	}
}

func TestTasksService_List(t *testing.T) {
	expected := TasksResponse{
		Data: []Task{
			fakeTask,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tasks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[action]":                 []string{"action"},
			"filter[status]":                 []string{"status"},
			"filter[compute_resource_id]":    []string{"1"},
			"filter[compute_resource_vm_id]": []string{"2"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterTasks{}).
		ByAction("action").
		ByStatus("status").
		ByComputeResourceID(1).
		ByComputeResourceVmID(2)

	actual, err := createTestClient(t, s.URL).Tasks.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestTasksService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tasks/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeTask)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Tasks.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeTask, actual)
}
