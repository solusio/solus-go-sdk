package solus

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/guregu/null.v4"
)

type LocationsService service

type LocationCreateRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	IconId           null.Int `json:"icon_id"`
	IsDefault        bool     `json:"is_default"`
	IsVisible        bool     `json:"is_visible"`
	ComputeResources []int    `json:"compute_resources"`
}

type Location struct {
	Id               int               `json:"id"`
	Name             string            `json:"name"`
	Icon             Icon              `json:"icon"`
	Description      string            `json:"description"`
	IsDefault        bool              `json:"is_default"`
	IsVisible        bool              `json:"is_visible"`
	ComputeResources []ComputeResource `json:"compute_resources"`
}

type LocationResponse struct {
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

	var resp LocationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Location{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

type LocationsPaginatedResponse struct {
	Data  []Location    `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

func (s *LocationsService) List(ctx context.Context, filter *FilterLocations) (LocationsPaginatedResponse, error) {
	opts := newRequestOpts()
	opts.params = filterToParams(filter.get())

	body, code, err := s.client.request(ctx, "GET", "locations", withParams(opts))
	if err != nil {
		return LocationsPaginatedResponse{}, err
	}

	if code != 200 {
		return LocationsPaginatedResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp LocationsPaginatedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return LocationsPaginatedResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}

func (s *LocationsService) Get(ctx context.Context, id int) (Location, error) {
	body, code, err := s.client.request(ctx, "GET", fmt.Sprintf("locations/%d", id))
	if err != nil {
		return Location{}, err
	}

	if code != 200 {
		return Location{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp LocationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Location{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *LocationsService) Update(ctx context.Context, id int, data LocationCreateRequest) (Location, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "PUT", fmt.Sprintf("locations/%d", id), withBody(opts))
	if err != nil {
		return Location{}, err
	}

	if code != 200 {
		return Location{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp LocationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Location{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *LocationsService) Delete(ctx context.Context, id int) error {
	body, code, err := s.client.request(ctx, "DELETE", fmt.Sprintf("locations/%d", id))
	if err != nil {
		return err
	}

	if code != 204 {
		return fmt.Errorf("HTTP %d: %s", code, body)
	}
	return nil
}
