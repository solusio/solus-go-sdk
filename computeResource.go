package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	ComputeResourceStatusActive           = "active"
	ComputeResourceStatusCommissioning    = "commissioning"
	ComputeResourceStatusConfigureNetwork = "configure_network"
	ComputeResourceStatusFailed           = "failed"
	ComputeResourceStatusUnavailable      = "unavailable"
)

type ComputeResource struct {
	Id                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	CanRetryInstallation bool                   `json:"can_retry_installation"`
	Host                 string                 `json:"host"`
	AgentPort            int                    `json:"agent_port"`
	Status               ComputerResourceStatus `json:"status"`
	Zones                []Zone                 `json:"zones"`
}

type ComputerResourceStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ComputerResourceCreateRequest struct {
	Name  string `json:"name"`
	Host  string `json:"host"`
	Login string `json:"login"`
	// SSH port number
	Port int `json:"port"`
	// Auth type 'lpass' or 'lkey'
	Type     string `json:"type"`
	Password string `json:"password"`
	// SSH private key
	Key       string `json:"key"`
	AgentPort int    `json:"agent_port"`
}

type ComputerResourceResponse struct {
	Data ComputeResource `json:"data"`
}

func (c *Client) ComputerResourceCreate(ctx context.Context, data ComputerResourceCreateRequest) (ComputeResource, error) {
	body, code, err := c.request(ctx, "POST", "compute_resources", data)
	if err != nil {
		return ComputeResource{}, err
	}

	if code != 201 {
		return ComputeResource{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ComputerResourceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ComputeResource{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) ComputerResource(ctx context.Context, id int) (ComputeResource, error) {
	body, code, err := c.request(ctx, "GET", fmt.Sprintf("compute_resources/%d", id), nil)
	if err != nil {
		return ComputeResource{}, err
	}

	if code != 200 {
		return ComputeResource{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ComputerResourceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ComputeResource{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
