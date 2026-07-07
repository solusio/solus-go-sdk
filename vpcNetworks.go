package solus

import (
	"context"
	"fmt"

	"gopkg.in/guregu/null.v4"
)

// VPCNetworksService handles all available methods with VPC networks.
type VPCNetworksService service

// VPCNetworkListType represents available list types for VPC networks.
type VPCNetworkListType string

const (
	// VPCNetworkListTypeRange indicates a range-based IP list.
	VPCNetworkListTypeRange VPCNetworkListType = "range"

	// VPCNetworkListTypeSet indicates a set-based IP list.
	VPCNetworkListTypeSet VPCNetworkListType = "set"
)

// VPCNetwork represents a VPC network.
type VPCNetwork struct {
	ID       int                `json:"id"`
	Name     string             `json:"name"`
	ListType VPCNetworkListType `json:"list_type"`
	From     string             `json:"from"`
	To       string             `json:"to"`
	Netmask  string             `json:"netmask"`
	Location ShortLocation      `json:"location"`
	User     User               `json:"user"`
}

// ShortVPCNetwork represents only ID and name of VPC network.
type ShortVPCNetwork struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// VPCNetworkCreateRequest represents available properties for creating a new VPC network.
type VPCNetworkCreateRequest struct {
	Name       string             `json:"name"`
	ListType   VPCNetworkListType `json:"list_type"`
	From       string             `json:"from,omitempty"`
	To         string             `json:"to,omitempty"`
	Netmask    string             `json:"netmask"`
	LocationID int                `json:"location_id"`
	UserID     null.Int           `json:"user_id,omitempty"`
}

// VPCNetworkUpdateRequest represents available properties for updating an existing VPC network.
type VPCNetworkUpdateRequest struct {
	Name     string             `json:"name"`
	ListType VPCNetworkListType `json:"list_type"`
	From     string             `json:"from,omitempty"`
	To       string             `json:"to,omitempty"`
	Netmask  string             `json:"netmask"`
	UserID   null.Int           `json:"user_id,omitempty"`
}

// VPCNetworkIP represents an IP address in a VPC network.
type VPCNetworkIP struct {
	ID         int             `json:"id"`
	IP         string          `json:"ip"`
	Server     VirtualServer   `json:"server"`
	VPCNetwork ShortVPCNetwork `json:"vpc_network"`
}

// VPCNetworkIPAddRequest represents a request to add a single IP to VPC network.
type VPCNetworkIPAddRequest struct {
	IP string `json:"ip"`
}

// VPCNetworkIPsAddRequest represents a request to add multiple IPs to VPC network.
type VPCNetworkIPsAddRequest struct {
	IPs []string `json:"ips"`
}

// VPCNetworkAttachRequest represents a request to attach VPC network to servers.
type VPCNetworkAttachRequest struct {
	ServerIDs []int `json:"servers_ids"`
}

// VPCNetworkDetachRequest represents a request to detach VPC network from servers.
type VPCNetworkDetachRequest struct {
	ServerIDs []int `json:"servers_ids"`
}

// VPCNetworksResponse represents paginated list of VPC networks.
// This cursor can be used for iterating over all available VPC networks.
type VPCNetworksResponse struct {
	paginatedResponse

	Data []VPCNetwork `json:"data"`
}

// VPCNetworkIPsResponse represents paginated list of VPC network IPs.
type VPCNetworkIPsResponse struct {
	paginatedResponse

	Data []VPCNetworkIP `json:"data"`
}

type vpcNetworkResponse struct {
	Data VPCNetwork `json:"data"`
}

// Create creates new VPC network.
func (s *VPCNetworksService) Create(ctx context.Context, data VPCNetworkCreateRequest) (VPCNetwork, error) {
	var resp vpcNetworkResponse
	return resp.Data, s.client.create(ctx, "vpc_networks", data, &resp)
}

// List lists VPC networks.
func (s *VPCNetworksService) List(ctx context.Context, filter *FilterVPCNetworks) (VPCNetworksResponse, error) {
	resp := VPCNetworksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "vpc_networks", &resp, withFilter(filter.data))
}

// Get gets specified VPC network.
func (s *VPCNetworksService) Get(ctx context.Context, id int) (VPCNetwork, error) {
	var resp vpcNetworkResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("vpc_networks/%d", id), &resp)
}

// Update updates specified VPC network.
func (s *VPCNetworksService) Update(ctx context.Context, id int, data VPCNetworkUpdateRequest) (VPCNetwork, error) {
	var resp vpcNetworkResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("vpc_networks/%d", id), data, &resp)
}

// Delete deletes specified VPC network.
func (s *VPCNetworksService) Delete(ctx context.Context, id int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("vpc_networks/%d", id))
}

// AddIP adds a single IP address to VPC network.
func (s *VPCNetworksService) AddIP(ctx context.Context, id int, data VPCNetworkIPAddRequest) (VPCNetworkIP, error) {
	var resp struct {
		Data VPCNetworkIP `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("vpc_networks/%d/add_ip", id), data, &resp)
}

// AddIPs adds multiple IP addresses to VPC network.
func (s *VPCNetworksService) AddIPs(ctx context.Context, id int, data VPCNetworkIPsAddRequest) error {
	return s.client.create(ctx, fmt.Sprintf("vpc_networks/%d/add_ips", id), data, nil)
}

// Attach attaches VPC network to servers.
func (s *VPCNetworksService) Attach(ctx context.Context, id int, data VPCNetworkAttachRequest) error {
	return s.client.post(ctx, fmt.Sprintf("vpc_networks/%d/attach", id), data, nil)
}

// Detach detaches VPC network from servers.
func (s *VPCNetworksService) Detach(ctx context.Context, id int, data VPCNetworkDetachRequest) error {
	return s.client.post(ctx, fmt.Sprintf("vpc_networks/%d/detach", id), data, nil)
}

// ListIPs lists IP addresses in VPC network.
func (s *VPCNetworksService) ListIPs(ctx context.Context, id int) (VPCNetworkIPsResponse, error) {
	resp := VPCNetworkIPsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, fmt.Sprintf("vpc_networks/%d/ips", id), &resp)
}

// DeleteIP deletes specified IP address from VPC network.
func (s *VPCNetworksService) DeleteIP(ctx context.Context, vpcNetworkID, ipID int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("vpc_networks/%d/ips/%d", vpcNetworkID, ipID))
}
