package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			DiskBandwidth: DiskBandwidthPlanLimit{
				IsEnabled: true,
				Limit:     11,
				Unit:      DiskBandwidthPlanLimitUnitBps,
			},
			DiskIOPS: DiskIOPSPlanLimit{
				IsEnabled: true,
				Limit:     12,
				Unit:      DiskIOPSPlanLimitUnitIOPS,
			},
			NetworkIncomingBandwidth: BandwidthPlanLimit{
				IsEnabled: true,
				Limit:     13,
				Unit:      BandwidthPlanLimitUnitKbps,
			},
			NetworkOutgoingBandwidth: BandwidthPlanLimit{
				IsEnabled: true,
				Limit:     14,
				Unit:      BandwidthPlanLimitUnitMbps,
			},
			NetworkIncomingTraffic: TrafficPlanLimit{
				IsEnabled: true,
				Limit:     15,
				Unit:      TrafficPlanLimitUnitTB,
			},
			NetworkOutgoingTraffic: TrafficPlanLimit{
				IsEnabled: true,
				Limit:     16,
				Unit:      TrafficPlanLimitUnitMB,
			},
			NetworkReduceBandwidth: BandwidthPlanLimit{},
		},
		TokensPerHour:    4,
		TokensPerMonth:   5,
		Position:         6,
		ResetLimitPolicy: PlanResetLimitPolicyVMCreatedDay,
	}

	expectedData := data
	expectedData.Limits.NetworkReduceBandwidth.Unit = BandwidthPlanLimitUnitKbps

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, expectedData)

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

func TestPlansService_setDefaultForPlanLimits(t *testing.T) {
	testCases := map[string]struct {
		given    PlanLimits
		expected PlanLimits
	}{
		"empty": {
			PlanLimits{},
			PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
			},
		},

		"with units": {
			PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
			},
			PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
			},
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			new(PlansService).setDefaultsForPlanLimits(&tt.given)
			assert.Equal(t, tt.expected, tt.given)
		})
	}
}
