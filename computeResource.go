package solus

import (
	"context"
	"fmt"
	"net/http"
)

type ComputeResourcesService service

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
	IpBlocks  []int  `json:"ip_blocks,omitempty"`
	Locations []int  `json:"locations,omitempty"`
}

type ComputerResourceResponse struct {
	Data ComputeResource `json:"data"`
}

type ComputerResourceNetworksResponse struct {
	Data []ComputeResourceNetwork `json:"data"`
}

func (s *ComputeResourcesService) Create(ctx context.Context, data ComputerResourceCreateRequest) (ComputeResource, error) {
	var resp ComputerResourceResponse
	return resp.Data, s.client.create(ctx, "compute_resources", data, &resp)
}

func (s *ComputeResourcesService) Get(ctx context.Context, id int) (ComputeResource, error) {
	var resp ComputerResourceResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d", id), &resp)
}

func (s *ComputeResourcesService) Networks(ctx context.Context, id int) ([]ComputeResourceNetwork, error) {
	var resp ComputerResourceNetworksResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("compute_resources/%d/networks", id), &resp)
}

func (s *ComputeResourcesService) SetUpNetwork(ctx context.Context, id int, networkId string) error {
	data := struct {
		Id string `json:"id"`
	}{
		Id: networkId,
	}
	body, code, err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("compute_resources/%d/setup_network", id), withBody(data))
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return newHTTPError(code, body)
	}
	return nil
}
