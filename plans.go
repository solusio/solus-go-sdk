package solus

import (
	"context"
	"encoding/json"
	"fmt"
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

	body, code, err := s.client.request(ctx, "GET", "plans")
	if err != nil {
		return PlansResponse{}, err
	}

	if code != 200 {
		return PlansResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return PlansResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}

func (s *PlansService) Create(ctx context.Context, data PlanCreateRequest) (Plan, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "plans", withBody(opts))
	if err != nil {
		return Plan{}, err
	}

	if code != 201 {
		return Plan{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp PlanCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Plan{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
