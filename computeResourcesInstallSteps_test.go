package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeResourcesService_InstallSteps(t *testing.T) {
	expected := []ComputeResourceInstallStep{
		fakeComputeResourceInstallStep,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/install_steps", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.InstallSteps(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestComputeResourcesService_InstallStepCreate(t *testing.T) {
	data := ComputeResourceInstallStepRequest{
		Title:      "fake title",
		Status:     ComputeResourceInstallStepStatusRunning,
		StatusText: "fake status text",
		Progress:   42,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/install_steps", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeComputeResourceInstallStep)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.InstallStepCreate(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResourceInstallStep, actual)
}

func TestComputeResourcesService_InstallStepUpdate(t *testing.T) {
	data := ComputeResourceInstallStepRequest{
		Title:      "fake title",
		Status:     ComputeResourceInstallStepStatusDone,
		StatusText: "done",
		Progress:   100,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/10/install_steps/5", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeComputeResourceInstallStep)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.InstallStepUpdate(context.Background(), 10, 5, data)
	require.NoError(t, err)
	require.Equal(t, fakeComputeResourceInstallStep, actual)
}
