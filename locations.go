package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type LocationsService service

type LocationCreateRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Icon             string `json:"icon"`
	IsDefault        bool   `json:"is_default"`
	IsVisible        bool   `json:"is_visible"`
	ComputeResources []int  `json:"compute_resources"`
}

type Location struct {
	Id               int               `json:"id"`
	Name             string            `json:"name"`
	Icon             string            `json:"icon"`
	Description      string            `json:"description"`
	IsDefault        bool              `json:"is_default"`
	IsVisible        bool              `json:"is_visible"`
	ComputeResources []ComputeResource `json:"compute_resources"`
}

type LocationCreateResponse struct {
	Data Location `json:"data"`
}

func (s *LocationsService) Create(ctx context.Context, data LocationCreateRequest) (Location, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "locations", withBody(opts))
	if err != nil {
		return Location{}, err
	}

	if code != 201 {
		return Location{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp LocationCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Location{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
