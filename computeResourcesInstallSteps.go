package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ComputeResourceInstallStepStatus string

const (
	ComputeResourceInstallStepStatusError ComputeResourceInstallStepStatus = "error"
)

type ComputeResourceInstallStepsResponse struct {
	Data []ComputeResourceInstallStep `json:"data"`
}

type ComputeResourceInstallStep struct {
	Id                int                              `json:"id"`
	ComputeResourceId int                              `json:"compute_resource_id"`
	Title             string                           `json:"title"`
	Status            ComputeResourceInstallStepStatus `json:"status"`
	StatusText        string                           `json:"status_text"`
	Progress          int                              `json:"progress"`
}

func (s *ComputeResourcesService) InstallSteps(ctx context.Context, id int) ([]ComputeResourceInstallStep, error) {
	body, code, err := s.client.request(ctx, "GET", fmt.Sprintf("compute_resources/%d/install_steps", id))
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ComputeResourceInstallStepsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
