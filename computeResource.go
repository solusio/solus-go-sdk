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
	Id        int                    `json:"id"`
	Name      string                 `json:"name"`
	Host      string                 `json:"host"`
	AgentPort int                    `json:"agent_port"`
	Status    ComputerResourceStatus `json:"status"`
	Locations []Location             `json:"locations"`
}

type ComputeResourceNetwork struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// AddrConfType for 'static' or 'dhcp'
	AddrConfType string `json:"addr_conf_type"`
	IpVersion    int    `json:"ip_version"`
	Ip           string `json:"ip"`
	Mask         string `json:"mask"`
	MaskSize     int    `json:"mask_size"`
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

type ComputerResourceNetworksResponse struct {
	Data []ComputeResourceNetwork `json:"data"`
}

func (c *Client) ComputerResourceCreate(ctx context.Context, data ComputerResourceCreateRequest) (ComputeResource, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := c.request(ctx, "POST", "compute_resources", withBody(opts))
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
	body, code, err := c.request(ctx, "GET", fmt.Sprintf("compute_resources/%d", id))
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

func (c *Client) ComputerResourceNetworks(ctx context.Context, id int) ([]ComputeResourceNetwork, error) {
	body, code, err := c.request(ctx, "GET", fmt.Sprintf("compute_resources/%d/networks", id))
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ComputerResourceNetworksResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) ComputerResourceSetUpNetwork(ctx context.Context, id int, networkId string) error {
	data := struct {
		Id string `json:"id"`
	}{
		Id: networkId,
	}
	opts := newRequestOpts()
	opts.body = data
	body, code, err := c.request(ctx, "POST", fmt.Sprintf("compute_resources/%d/setup_network", id), withBody(opts))
	if err != nil {
		return err
	}

	if code != 200 {
		return fmt.Errorf("HTTP %d: %s", code, body)
	}

	return nil
}
