package solus

import (
	"context"
)

type PlansService service

type PlanParams struct {
	Disk int `json:"disk"`
	RAM  int `json:"ram"`
	VCPU int `json:"vcpu"`
}

type PlanLimit struct {
	IsEnabled bool `json:"is_enabled"`
	Limit     int  `json:"limit"`
}

type PlanLimits struct {
	TotalBytes PlanLimit `json:"total_bytes"`
	TotalIops  PlanLimit `json:"total_iops"`
}

type Plan struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Params              PlanParams `json:"params"`
	StorageType         string     `json:"storage_type"`
	ImageFormat         string     `json:"image_format"`
	IsDefault           bool       `json:"is_default"`
	IsSnapshotAvailable bool       `json:"is_snapshot_available"`
	IsSnapshotsEnabled  bool       `json:"is_snapshots_enabled"`
	Limits              PlanLimits `json:"limits"`
	TokenPerHour        float64    `json:"token_per_hour"`
	TokenPerMonth       float64    `json:"token_per_month"`
	Position            float64    `json:"position"`
}

type PlanCreateRequest struct {
	Name               string     `json:"name"`
	Type               string     `json:"type"`
	Params             PlanParams `json:"params"`
	StorageType        string     `json:"storage_type"`
	ImageFormat        string     `json:"image_format"`
	IsVisible          bool       `json:"is_visible"`
	IsSnapshotsEnabled bool       `json:"is_snapshots_enabled"`
	Limits             PlanLimits `json:"limits"`
	TokenPerHour       float64    `json:"token_per_hour"`
	TokenPerMonth      float64    `json:"token_per_month"`
	Position           float64    `json:"position"`
}

type PlansResponse struct {
	paginatedResponse

	Data []Plan `json:"data"`
}

type PlanCreateResponse struct {
	Data Plan `json:"data"`
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
	var resp PlanCreateResponse
	return resp.Data, s.client.create(ctx, "plans", data, &resp)
}
