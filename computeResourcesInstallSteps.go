package solus

import (
	"context"
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
	var resp ComputeResourceInstallStepsResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/install_steps", id), &resp)
}
