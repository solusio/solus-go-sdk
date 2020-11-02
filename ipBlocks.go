package solus

import (
	"context"
	"fmt"
	"net/http"
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
	return resp, s.client.list(ctx, "ip_blocks", &resp)
}

func (s *IpBlocksService) Create(ctx context.Context, data IpBlockCreateRequest) (IpBlock, error) {
	var resp IpBlockCreateResponse
	return resp.Data, s.client.create(ctx, "ip_blocks", data, &resp)
}

func (s *IpBlocksService) IpAddressCreate(ctx context.Context, ipBlockId int) (IpBlockIpAddress, error) {
	body, code, err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("ip_blocks/%d/ips", ipBlockId))
	if err != nil {
		return IpBlockIpAddress{}, err
	}

	if code != http.StatusCreated {
		return IpBlockIpAddress{}, newHTTPError(code, body)
	}

	var resp IpBlockIpAddressCreateResponse
	return resp.Data, unmarshal(body, &resp)
}

func (s *IpBlocksService) IpAddressDelete(ctx context.Context, ipId int) error {
	return s.client.delete(ctx, fmt.Sprintf("ips/%d", ipId))
}
