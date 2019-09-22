package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type Plan struct {
	Id                  int        `json:"id"`
	Name                string     `json:"name"`
	Type                string     `json:"type"`
	Params              PlanParams `json:"params"`
	StorageType         string     `json:"storage_type"`
	IsDefault           bool       `json:"is_default"`
	IsSnapshotAvailable bool       `json:"is_snapshot_available"`
	IsSnapshotEnabled   bool       `json:"is_snapshot_enabled"`
	TokenValue          int        `json:"token_value"`
}

type PlanParams struct {
	Hdd int `json:"hdd"`
	Ram int `json:"ram"`
	Cpu int `json:"cpu"`
}

type PlansResponse struct {
	Data  []Plan        `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

func (c *Client) Plans(ctx context.Context) ([]Plan, error) {
	body, code, err := c.request(ctx, "POST", "plans", nil)
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
