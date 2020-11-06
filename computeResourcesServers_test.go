package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestComputeResourcesService_ServersCreate(t *testing.T) {
	data := ComputeResourceServerCreateRequest{
		Name:             "fake name",
		Description:      "fake description",
		Password:         "123456789",
		PlanID:           1,
		OSImageVersionID: 2,
		ApplicationID:    3,
		ApplicationData: map[string]string{
			"foo": "bar",
		},
		SSHKeys:   []int{4, 5},
		UserData:  "fake user data",
		FQDNs:     []string{"example.com"},
		UserID:    6,
		ProjectID: 7,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/compute_resources/42/servers", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeServer)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).ComputeResources.ServersCreate(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeServer, actual)
}
