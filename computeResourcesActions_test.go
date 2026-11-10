package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResourcesService_IPBlocksList(t *testing.T) {
	expected := []IPBlock{fakeIPBlock}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/ip_blocks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.IPBlocksList(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_Metrics(t *testing.T) {
	expected := ComputeResourceMetric{
		Total:             10,
		User:              4,
		Name:              "cpu",
		ComputeResourceID: 10,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/metrics", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Metrics(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_TasksList(t *testing.T) {
	expected := []Task{fakeTask}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/tasks", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.TasksList(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_Usage(t *testing.T) {
	expected := ComputeResourceUsage{
		ServerCount: 3,
		CPU:         4,
		RAM:         5,
		Disk:        6,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/usage", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.Usage(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_RetryInstallation(t *testing.T) {
	data := ComputerResourceCreateRequest{
		Name:      "name",
		Host:      "host",
		Login:     "login",
		Port:      22,
		Type:      ComputeResourceAuthTypePassword,
		Password:  "password",
		AgentPort: 1234,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/retry_installation", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeComputeResource)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.RetryInstallation(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResource, actual)
}
