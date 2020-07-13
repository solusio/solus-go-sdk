package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type IpBlocksService service

type NetworkType string

const (
	Routed  NetworkType = "routed"
	Bridged NetworkType = "bridged"
)

type IpBlockCreateRequest struct {
	Name             string      `json:"name"`
	NetworkType      NetworkType `json:"network_type"`
	Gateway          string      `json:"gateway"`
	Netmask          string      `json:"netmask"`
	Ns1              string      `json:"ns_1"`
	Ns2              string      `json:"ns_2"`
	ComputeResources []int       `json:"compute_resources"`
	From             string      `json:"from"`
	To               string      `json:"to"`
}

type IpBlock struct {
	Id               int                `json:"id"`
	Name             string             `json:"name"`
	Gateway          string             `json:"gateway"`
	Netmask          string             `json:"netmask"`
	Ns1              string             `json:"ns_1"`
	Ns2              string             `json:"ns_2"`
	From             string             `json:"from"`
	To               string             `json:"to"`
	ComputeResources []ComputeResource  `json:"compute_resources[]"`
	Ips              []IpBlockIpAddress `json:"ips[]"`
	NetworkType      NetworkType        `json:"network_type"`
}

type IpBlockCreateResponse struct {
	Data IpBlock `json:"data"`
}

type IpBlockIpAddress struct {
	Id int    `json:"id"`
	Ip string `json:"ip"`
}

type IpBlockIpAddressCreateResponse struct {
	Data IpBlockIpAddress `json:"data"`
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
