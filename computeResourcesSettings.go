package solus

import (
	"context"
	"fmt"
)

// SettingsGet gets compute resource's settings.
func (s *ComputeResourcesService) SettingsGet(ctx context.Context, id int) (ComputeResourceSettings, error) {
	var resp struct {
		Data ComputeResourceSettings `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/settings", id), &resp)
}

// SettingsUpdate updates compute resource's settings.
func (s *ComputeResourcesService) SettingsUpdate(
	ctx context.Context,
	id int,
	data ComputeResourceSettings,
) (ComputeResourceSettings, error) {
	var resp struct {
		Data ComputeResourceSettings `json:"data"`
	}
	return resp.Data, s.client.update(ctx, fmt.Sprintf("compute_resources/%d/settings", id), data, &resp)
}
