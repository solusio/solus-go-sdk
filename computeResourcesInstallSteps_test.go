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
