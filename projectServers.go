package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ProjectServersCreateRequest struct {
	Name             string `json:"name"`
	PlanId           int    `json:"plan_id"`
	LocationId       int    `json:"location_id"`
	OsImageVersionId int    `json:"os_image_version_id"`
	SshKeys          []int  `json:"ssh_keys,omitempty"`
	UserData         string `json:"user_data,omitempty"`
}

type ProjectServersCreateResponse struct {
	Data Server `json:"data"`
}

type ProjectServersResponse struct {
	paginatedResponse

	Data []Server `json:"data"`
}

type Server struct {
	Id          int                `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	UUID        string             `json:"uuid"`
	Status      string             `json:"status"`
	Ips         []IpBlockIpAddress `json:"ips"`
}

func (s *ProjectsService) ServersCreate(ctx context.Context, projectId int, data ProjectServersCreateRequest) (Server, error) {
	body, code, err := s.client.request(ctx, "POST", fmt.Sprintf("projects/%d/servers", projectId), withBody(data))
	if err != nil {
		return Server{}, err
	}

	if code != 201 {
		return Server{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ProjectServersCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Server{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *ProjectsService) ServersListAll(ctx context.Context, projectId int) ([]Server, error) {
	resp, err := s.Servers(ctx, projectId)
	if err != nil {
		return nil, err
	}

	servers := make([]Server, len(resp.Data))
	copy(servers, resp.Data)
	for resp.Next(ctx) {
		servers = append(servers, resp.Data...)
	}
	return servers, resp.Err()
}

func (s *ProjectsService) Servers(ctx context.Context, projectId int) (ProjectServersResponse, error) {
	resp := ProjectServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	body, code, err := s.client.request(ctx, "GET", fmt.Sprintf("projects/%d/servers", projectId))
	if err != nil {
		return ProjectServersResponse{}, err
	}

	if code != 200 {
		return ProjectServersResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return ProjectServersResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}
