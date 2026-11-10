package solus

import (
	"context"
	"fmt"
	"net/http"
)

// ComputeResourceMetric represents a metric value for a compute resource.
type ComputeResourceMetric struct {
	Total             float64 `json:"total"`
	User              float64 `json:"user"`
	Name              string  `json:"name"`
	ComputeResourceID int     `json:"compute_resource_id"`
}

// ComputeResourceUsage represents resource usage by a compute resource.
type ComputeResourceUsage struct {
	ServerCount int `json:"server_count"`
	CPU         int `json:"cpu"`
	RAM         int `json:"ram"`
	Disk        int `json:"disk"`
}

// IPBlocksList lists IP blocks assigned to the specified compute resource.
func (s *ComputeResourcesService) IPBlocksList(ctx context.Context, id int) ([]IPBlock, error) {
	var resp struct {
		Data []IPBlock `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/ip_blocks", id), &resp)
}

// Metrics gets metrics for the specified compute resource.
func (s *ComputeResourcesService) Metrics(ctx context.Context, id int) (ComputeResourceMetric, error) {
	var resp struct {
		Data ComputeResourceMetric `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/metrics", id), &resp)
}

// TasksList lists tasks related to the specified compute resource.
func (s *ComputeResourcesService) TasksList(ctx context.Context, id int) ([]Task, error) {
	var resp struct {
		Data []Task `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/tasks", id), &resp)
}

// Usage gets resource usage for the specified compute resource.
func (s *ComputeResourcesService) Usage(ctx context.Context, id int) (ComputeResourceUsage, error) {
	var resp struct {
		Data ComputeResourceUsage `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/usage", id), &resp)
}

// RetryInstallation retries installation of the specified compute resource.
func (s *ComputeResourcesService) RetryInstallation(
	ctx context.Context,
	id int,
	data ComputerResourceCreateRequest,
) (ComputeResource, error) {
	path := fmt.Sprintf("compute_resources/%d/retry_installation", id)
	body, code, err := s.client.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return ComputeResource{}, err
	}

	if code != http.StatusOK {
		return ComputeResource{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp computeResourceResponse
	return resp.Data, unmarshal(body, &resp)
}
