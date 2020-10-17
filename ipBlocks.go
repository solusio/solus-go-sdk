package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type IpBlocksService service

type IPVersion string

const (
	IPv4 IPVersion = "IPv4"
	IPv6 IPVersion = "IPv6"
)

type IpBlockCreateRequest struct {
	ComputeResources []int     `json:"compute_resources,omitempty"`
	Name             string    `json:"name"`
	Type             IPVersion `json:"type"`
	Gateway          string    `json:"gateway"`
	Ns1              string    `json:"ns_1"`
	Ns2              string    `json:"ns_2"`

	// IPv4 related fields
	Netmask string `json:"netmask"`
	From    string `json:"from"`
	To      string `json:"to"`

	// IPv6 related fields
	Range  string `json:"range"`
	Subnet int    `json:"subnet"`
}

type IpBlock struct {
	Id               int                `json:"id"`
	Name             string             `json:"name"`
	Type             IPVersion          `json:"type"`
	Gateway          string             `json:"gateway"`
	Netmask          string             `json:"netmask"`
	Ns1              string             `json:"ns_1"`
	Ns2              string             `json:"ns_2"`
	From             string             `json:"from"`
	To               string             `json:"to"`
	Subnet           int                `json:"subnet"`
	ComputeResources []ComputeResource  `json:"compute_resources[]"`
	Ips              []IpBlockIpAddress `json:"ips[]"`
}

type IpBlockCreateResponse struct {
	Data IpBlock `json:"data"`
}

type IpBlockIpAddress struct {
	Id      int     `json:"id"`
	Ip      string  `json:"ip"`
	IpBlock IpBlock `json:"ip_block"`
}

type IpBlockIpAddressCreateResponse struct {
	Data IpBlockIpAddress `json:"data"`
}

type IpBlocksResponse struct {
	paginatedResponse

	Data []IpBlock `json:"data"`
}

func (s *IpBlocksService) List(ctx context.Context) (IpBlocksResponse, error) {
	resp := IpBlocksResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}

	body, code, err := s.client.request(ctx, "GET", "ip_blocks")
	if err != nil {
		return IpBlocksResponse{}, err
	}

	if code != 200 {
		return IpBlocksResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return IpBlocksResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}

func (s *IpBlocksService) Create(ctx context.Context, data IpBlockCreateRequest) (IpBlock, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "ip_blocks", withBody(opts))
	if err != nil {
		return IpBlock{}, err
	}

	if code != 201 {
		return IpBlock{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp IpBlockCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return IpBlock{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *IpBlocksService) IpAddressCreate(ctx context.Context, ipBlockId int) (IpBlockIpAddress, error) {
	body, code, err := s.client.request(ctx, "POST", fmt.Sprintf("ip_blocks/%d/ips", ipBlockId))
	if err != nil {
		return IpBlockIpAddress{}, err
	}

	if code != 201 {
		return IpBlockIpAddress{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp IpBlockIpAddressCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return IpBlockIpAddress{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *IpBlocksService) IpAddressDelete(ctx context.Context, ipId int) error {
	body, code, err := s.client.request(ctx, "DELETE", fmt.Sprintf("ips/%d", ipId))
	if err != nil {
		return err
	}

	if code != 204 {
		return fmt.Errorf("HTTP %d: %s", code, body)
	}

	return nil
}
