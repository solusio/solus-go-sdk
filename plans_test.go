package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskBandwidthPlanLimit_setDefault(t *testing.T) {
	testCases := map[string]*DiskBandwidthPlanLimit{
		"empty":     {},
		"with unit": {Unit: DiskBandwidthPlanLimitUnitBps},
	}

	for name, s := range testCases {
		t.Run(name, func(t *testing.T) {
			s.setDefault()
			assert.Equal(t, DiskBandwidthPlanLimitUnitBps, s.Unit)
		})
	}
}

func TestBandwidthPlanLimit_setDefault(t *testing.T) {
	testCases := map[BandwidthPlanLimitUnit]*BandwidthPlanLimit{
		BandwidthPlanLimitUnitKbps: {},
		BandwidthPlanLimitUnitGbps: {Unit: BandwidthPlanLimitUnitGbps},
	}

	for expected, s := range testCases {
		t.Run(string(expected), func(t *testing.T) {
			s.setDefault()
			assert.Equal(t, expected, s.Unit)
		})
	}
}

func TestDiskIOPSPlanLimit_setDefault(t *testing.T) {
	testCases := map[string]*DiskIOPSPlanLimit{
		"empty":     {},
		"with unit": {Unit: DiskIOPSPlanLimitUnitIOPS},
	}

	for name, s := range testCases {
		t.Run(name, func(t *testing.T) {
			s.setDefault()
			assert.Equal(t, DiskIOPSPlanLimitUnitIOPS, s.Unit)
		})
	}
}

func TestTrafficPlanLimit_setDefault(t *testing.T) {
	testCases := map[TrafficPlanLimitUnit]*TrafficPlanLimit{
		TrafficPlanLimitUnitKiB: {},
		TrafficPlanLimitUnitPiB: {Unit: TrafficPlanLimitUnitPiB},
	}

	for expected, s := range testCases {
		t.Run(string(expected), func(t *testing.T) {
			s.setDefault()
			assert.Equal(t, expected, s.Unit)
		})
	}
}

func TestUnitPlanLimit_setDefault(t *testing.T) {
	testCases := map[string]*UnitPlanLimit{
		"empty":     {},
		"with unit": {Unit: PlanLimitUnits},
	}

	for name, s := range testCases {
		t.Run(name, func(t *testing.T) {
			s.setDefault()
			assert.Equal(t, PlanLimitUnits, s.Unit)
		})
	}
}

func TestPlansService_List(t *testing.T) {
	expected := PlansResponse{
		Data: []Plan{
			fakePlan,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterPlans{}).ByName("name")

	actual, err := createTestClient(t, s.URL).Plans.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestPlansService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakePlan)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Plans.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakePlan, actual)
}

func TestPlansService_Create(t *testing.T) {
	data := PlanCreateRequest{
		Name: "name",
		Params: PlanParams{
			Disk: 1,
			RAM:  2,
			VCPU: 3,
		},
		VirtualizationType: VirtualizationTypeKVM,
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
				Unit:      TrafficPlanLimitUnitTiB,
			},
			NetworkOutgoingTraffic: TrafficPlanLimit{
				IsEnabled: true,
				Limit:     16,
				Unit:      TrafficPlanLimitUnitMiB,
			},
			NetworkTotalTraffic: TrafficPlanLimit{
				IsEnabled: false,
				Limit:     17,
				Unit:      TrafficPlanLimitUnitMiB,
			},
			NetworkReduceBandwidth: BandwidthPlanLimit{},
			BackupsNumber: UnitPlanLimit{
				IsEnabled: true,
				Limit:     7,
				Unit:      PlanLimitUnits,
			},
		},
		TokensPerHour:  4,
		TokensPerMonth: 5,
		BackupSettings: PlanBackupSettings{
			IsIncrementalBackupEnabled: false,
			IncrementalBackupsLimit:    3,
		},
		ResetLimitPolicy:        PlanResetLimitPolicyVMCreatedDay,
		NetworkTotalTrafficType: PlanNetworkTotalTrafficTypeSeparate,
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

func TestPlansService_setCreateRequestDefaults(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		req := PlanCreateRequest{}
		(&PlansService{}).setCreateRequestDefaults(&req)
		assert.Equal(t, PlanCreateRequest{
			Limits: PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyNever,
		}, req)
	})

	t.Run("full", func(t *testing.T) {
		req := PlanCreateRequest{
			Limits: PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitPiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyFirstDayOfMonth,
		}
		(&PlansService{}).setCreateRequestDefaults(&req)
		assert.Equal(t, PlanCreateRequest{
			Limits: PlanLimits{
				DiskBandwidth:            DiskBandwidthPlanLimit{Unit: DiskBandwidthPlanLimitUnitBps},
				DiskIOPS:                 DiskIOPSPlanLimit{Unit: DiskIOPSPlanLimitUnitIOPS},
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitPiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyFirstDayOfMonth,
		}, req)
	})
}

func TestPlansService_Update(t *testing.T) {
	data := PlanUpdateRequest{
		Name:               "name",
		IsVisible:          true,
		IsSnapshotsEnabled: true,
		Limits: PlanUpdateLimits{
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
				Unit:      TrafficPlanLimitUnitTiB,
			},
			NetworkOutgoingTraffic: TrafficPlanLimit{
				IsEnabled: true,
				Limit:     16,
				Unit:      TrafficPlanLimitUnitMiB,
			},
			NetworkTotalTraffic: TrafficPlanLimit{
				IsEnabled: false,
				Limit:     17,
				Unit:      TrafficPlanLimitUnitMiB,
			},
			NetworkReduceBandwidth: BandwidthPlanLimit{},
			BackupsNumber: UnitPlanLimit{
				IsEnabled: true,
				Limit:     7,
				Unit:      PlanLimitUnits,
			},
		},
		TokensPerHour:  4,
		TokensPerMonth: 5,
		BackupSettings: PlanBackupSettings{
			IsIncrementalBackupEnabled: false,
			IncrementalBackupsLimit:    3,
		},
		ResetLimitPolicy:        PlanResetLimitPolicyVMCreatedDay,
		NetworkTotalTrafficType: PlanNetworkTotalTrafficTypeSeparate,
	}

	expectedData := data
	expectedData.Limits.NetworkReduceBandwidth.Unit = BandwidthPlanLimitUnitKbps

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, expectedData)

		writeResponse(t, w, http.StatusOK, fakePlan)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Plans.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakePlan, actual)
}

func TestPlansService_setUpdateRequestDefaults(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		req := PlanUpdateRequest{}
		(&PlansService{}).setUpdateRequestDefaults(&req)
		assert.Equal(t, PlanUpdateRequest{
			Limits: PlanUpdateLimits{
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitKiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyNever,
		}, req)
	})

	t.Run("full", func(t *testing.T) {
		req := PlanUpdateRequest{
			Limits: PlanUpdateLimits{
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitPiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyFirstDayOfMonth,
		}
		(&PlansService{}).setUpdateRequestDefaults(&req)
		assert.Equal(t, PlanUpdateRequest{
			Limits: PlanUpdateLimits{
				NetworkIncomingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitMbps},
				NetworkOutgoingBandwidth: BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitGbps},
				NetworkIncomingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitMiB},
				NetworkOutgoingTraffic:   TrafficPlanLimit{Unit: TrafficPlanLimitUnitTiB},
				NetworkTotalTraffic:      TrafficPlanLimit{Unit: TrafficPlanLimitUnitPiB},
				NetworkReduceBandwidth:   BandwidthPlanLimit{Unit: BandwidthPlanLimitUnitKbps},
				BackupsNumber:            UnitPlanLimit{Unit: PlanLimitUnits},
			},
			ResetLimitPolicy: PlanResetLimitPolicyFirstDayOfMonth,
		}, req)
	})
}

func TestPlansService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/plans/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Plans.Delete(context.Background(), 10)
	require.NoError(t, err)
}
