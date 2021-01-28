package solus

import (
	"context"
	"fmt"
)

type PlansService service

type PlanParams struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

type DiskBandwidthPlanLimit struct {
	IsEnabled bool                       `json:"is_enabled"`
	Limit     int                        `json:"limit"`
	Unit      DiskBandwidthPlanLimitUnit `json:"unit"`
}

type DiskBandwidthPlanLimitUnit string

const (
	DiskBandwidthPlanLimitUnitBps DiskBandwidthPlanLimitUnit = "Bps"
)

type BandwidthPlanLimit struct {
	IsEnabled bool                   `json:"is_enabled"`
	Limit     int                    `json:"limit"`
	Unit      BandwidthPlanLimitUnit `json:"unit"`
}

type BandwidthPlanLimitUnit string

const (
	BandwidthPlanLimitUnitKbps BandwidthPlanLimitUnit = "Kbps"
	BandwidthPlanLimitUnitMbps BandwidthPlanLimitUnit = "Mbps"
	BandwidthPlanLimitUnitGbps BandwidthPlanLimitUnit = "Gbps"
)

type DiskIOPSPlanLimit struct {
	IsEnabled bool                  `json:"is_enabled"`
	Limit     int                   `json:"limit"`
	Unit      DiskIOPSPlanLimitUnit `json:"unit"`
}

type DiskIOPSPlanLimitUnit string

const (
	DiskIOPSPlanLimitUnitOPS DiskIOPSPlanLimitUnit = "ops"
)

type TrafficPlanLimit struct {
	IsEnabled bool                 `json:"is_enabled"`
	Limit     int                  `json:"limit"`
	Unit      TrafficPlanLimitUnit `json:"unit"`
}

type TrafficPlanLimitUnit string

const (
	TrafficPlanLimitUnitKB TrafficPlanLimitUnit = "KB"
	TrafficPlanLimitUnitMB TrafficPlanLimitUnit = "MB"
	TrafficPlanLimitUnitGB TrafficPlanLimitUnit = "GB"
	TrafficPlanLimitUnitTB TrafficPlanLimitUnit = "TB"
	TrafficPlanLimitUnitPB TrafficPlanLimitUnit = "PB"
)

type PlanLimits struct {
	DiskBandwidth            DiskBandwidthPlanLimit `json:"disk_bandwidth"`
	DiskIOPS                 DiskIOPSPlanLimit      `json:"disk_iops"`
	NetworkIncomingBandwidth BandwidthPlanLimit     `json:"network_incoming_bandwidth"`
	NetworkOutgoingBandwidth BandwidthPlanLimit     `json:"network_outgoing_bandwidth"`
	NetworkIncomingTraffic   TrafficPlanLimit       `json:"network_incoming_traffic"`
	NetworkOutgoingTraffic   TrafficPlanLimit       `json:"network_outgoing_traffic"`
	NetworkReduceBandwidth   BandwidthPlanLimit     `json:"network_reduce_bandwidth"`
}

type PlanResetLimitPolicy string

const (
	PlanResetLimitPolicyFirstDayOfMonth PlanResetLimitPolicy = "first_day_of_month"
	PlanResetLimitPolicyVMCreatedDay    PlanResetLimitPolicy = "vm_created_day"
)

type PlanPrice struct {
	PerHour        string        `json:"per_hour"`
	PerMonth       string        `json:"per_month"`
	CurrencyCode   string        `json:"currency_code"`
	TaxesInclusive bool          `json:"taxes_inclusive"`
	Taxes          []interface{} `json:"taxes"`
	TotalPrice     string        `json:"total_price"`
	BackupPrice    string        `json:"backup_price"`
}

type Plan struct {
	ID                  int                  `json:"id"`
	Name                string               `json:"name"`
	Params              PlanParams           `json:"params"`
	StorageType         string               `json:"storage_type"`
	ImageFormat         string               `json:"image_format"`
	IsDefault           bool                 `json:"is_default"`
	IsSnapshotAvailable bool                 `json:"is_snapshot_available"`
	IsSnapshotsEnabled  bool                 `json:"is_snapshots_enabled"`
	IsBackupAvailable   bool                 `json:"is_backup_available"`
	BackupPrice         float32              `json:"backup_price"`
	IsVisible           bool                 `json:"is_visible"`
	Limits              PlanLimits           `json:"limits"`
	TokensPerHour       float64              `json:"tokens_per_hour"`
	TokensPerMonth      float64              `json:"tokens_per_month"`
	Position            float64              `json:"position"`
	Price               PlanPrice            `json:"price"`
	ResetLimitPolicy    PlanResetLimitPolicy `json:"reset_limit_policy"`
}

type PlanCreateRequest struct {
	Name               string               `json:"name"`
	Params             PlanParams           `json:"params"`
	StorageType        StorageTypeName      `json:"storage_type"`
	ImageFormat        ImageFormat          `json:"image_format"`
	Limits             PlanLimits           `json:"limits"`
	TokensPerHour      float64              `json:"tokens_per_hour"`
	TokensPerMonth     float64              `json:"tokens_per_month"`
	Position           float64              `json:"position"`
	IsVisible          bool                 `json:"is_visible"`
	IsDefault          bool                 `json:"is_default"`
	IsSnapshotsEnabled bool                 `json:"is_snapshots_enabled"`
	IsBackupAvailable  bool                 `json:"is_backup_available"`
	BackupPrice        float32              `json:"backup_price"`
	ResetLimitPolicy   PlanResetLimitPolicy `json:"reset_limit_policy"`
}

type PlanUpdateRequest struct {
	Name               string               `json:"name"`
	Limits             PlanLimits           `json:"limits"`
	TokensPerHour      float64              `json:"tokens_per_hour"`
	TokensPerMonth     float64              `json:"tokens_per_month"`
	Position           float64              `json:"position"`
	IsVisible          bool                 `json:"is_visible"`
	IsDefault          bool                 `json:"is_default"`
	IsSnapshotsEnabled bool                 `json:"is_snapshots_enabled"`
	IsBackupAvailable  bool                 `json:"is_backup_available"`
	BackupPrice        float32              `json:"backup_price"`
	ResetLimitPolicy   PlanResetLimitPolicy `json:"reset_limit_policy"`
}

type planResponse struct {
	Data Plan `json:"data"`
}

type PlansResponse struct {
	paginatedResponse

	Data []Plan `json:"data"`
}

func (s *PlansService) List(ctx context.Context) (PlansResponse, error) {
	resp := PlansResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "plans", &resp)
}

func (s *PlansService) Create(ctx context.Context, data PlanCreateRequest) (Plan, error) {
	var resp planResponse
	data.Limits = s.setDefaultForPlanLimits(data.Limits)
	return resp.Data, s.client.create(ctx, "plans", data, &resp)
}

func (s *PlansService) Update(ctx context.Context, id int, data PlanUpdateRequest) (Plan, error) {
	var resp planResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("plans/%d", id), data, &resp)
}

func (s *PlansService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("plans/%d", id))
}

func (*PlansService) setDefaultForPlanLimits(p PlanLimits) PlanLimits {
	if p.DiskBandwidth.Unit == "" {
		p.DiskBandwidth.Unit = DiskBandwidthPlanLimitUnitBps
	}

	if p.DiskIOPS.Unit == "" {
		p.DiskIOPS.Unit = DiskIOPSPlanLimitUnitOPS
	}

	if p.NetworkIncomingBandwidth.Unit == "" {
		p.NetworkIncomingBandwidth.Unit = BandwidthPlanLimitUnitKbps
	}

	if p.NetworkOutgoingBandwidth.Unit == "" {
		p.NetworkOutgoingBandwidth.Unit = BandwidthPlanLimitUnitKbps
	}

	if p.NetworkIncomingTraffic.Unit == "" {
		p.NetworkIncomingTraffic.Unit = TrafficPlanLimitUnitKB
	}

	if p.NetworkOutgoingTraffic.Unit == "" {
		p.NetworkOutgoingTraffic.Unit = TrafficPlanLimitUnitKB
	}

	if p.NetworkReduceBandwidth.Unit == "" {
		p.NetworkReduceBandwidth.Unit = BandwidthPlanLimitUnitKbps
	}
	return p
}
