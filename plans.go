package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type PlanParams struct {
	Hdd int `json:"hdd"`
	Ram int `json:"ram"`
	Cpu int `json:"cpu"`
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
	Type                string     `json:"type"`
	Params              PlanParams `json:"params"`
	StorageType         string     `json:"storage_type"`
	IsDefault           bool       `json:"is_default"`
	IsSnapshotAvailable bool       `json:"is_snapshot_available"`
	IsSnapshotEnabled   bool       `json:"is_snapshot_enabled"`
	Limits              PlanLimits `json:"limits"`
	TokenPerHour       float64     `json:"token_per_hour"`
	TokenPerMonth      float64     `json:"token_per_month"`
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
}

type PlansResponse struct {
	Data  []Plan        `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

type PlanCreateResponse struct {
	Data Plan `json:"data"`
}

func (c *Client) Plans(ctx context.Context) ([]Plan, error) {
	body, code, err := c.request(ctx, "GET", "plans", nil)
	if err != nil {
		return []Plan{}, err
	}

	if code != 200 {
		return []Plan{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp PlansResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return []Plan{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) PlanCreate(ctx context.Context, data PlanCreateRequest) (Plan, error) {
	body, code, err := c.request(ctx, "POST", "plans", data)
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
