package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ZoneCreateRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Icon             string `json:"icon"`
	IsDefault        bool   `json:"is_default"`
	IsVisible        bool   `json:"is_visible"`
	ComputeResources []int  `json:"compute_resources"`
}

type Zone struct {
	Id               int               `json:"id"`
	Name             string            `json:"name"`
	Icon             string            `json:"icon"`
	Description      string            `json:"description"`
	IsDefault        bool              `json:"is_default"`
	IsVisible        bool              `json:"is_visible"`
	ComputeResources []ComputeResource `json:"compute_resources"`
}

type ZoneCreateResponse struct {
	Data Zone `json:"data"`
}

func (c *Client) ZoneCreate(ctx context.Context, data ZoneCreateRequest) (Zone, error) {
	body, code, err := c.request(ctx, "POST", "zones", data)
	if err != nil {
		return Zone{}, err
	}

	if code != 201 {
		return Zone{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ZoneCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Zone{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
