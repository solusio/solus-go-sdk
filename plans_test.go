package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestPlansService_List(t *testing.T) {
	expected := PlansResponse{
		Data: []Plan{
			fakePlan,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Plans.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestPlansService_Create(t *testing.T) {
	data := PlanCreateRequest{
		Name: "name",
		Params: PlanParams{
			Disk: 1,
			RAM:  2,
			VCPU: 3,
		},
		StorageType:        "storage type",
		ImageFormat:        "image format",
		IsVisible:          true,
		IsSnapshotsEnabled: true,
		Limits: PlanLimits{
			TotalBytes: PlanLimit{
				IsEnabled: true,
				Limit:     4,
			},
			TotalIops: PlanLimit{
				IsEnabled: true,
				Limit:     5,
			},
		},
		TokensPerHour:  6,
		TokensPerMonth: 7,
		Position:       8,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakePlan)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Plans.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakePlan, actual)
}

func TestPlansService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Plans.Delete(context.Background(), 10)
	require.NoError(t, err)
}
