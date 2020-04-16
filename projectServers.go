package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ProjectServersCreateRequest struct {
	Name             string `json:"name"`
	PlanId           int    `json:"plan_id"`
	LocationId       int    `json:"zone_id"`
	OsImageVersionId int    `json:"os_image_version_id"`
	SshKeys          []int  `json:"ssh_keys,omitempty"`
	UserData         string `json:"user_data,omitempty"`
}

type ProjectServersCreateResponse struct {
	Data Server `json:"data"`
}

type ProjectServersResponse struct {
	Data  []Server      `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

type Server struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UUID        string `json:"uuid"`
	Status      string `json:"status"`
	Ip          string `json:"ip"`
}

func (c *Client) ProjectServerCreate(ctx context.Context, projectId int, data ProjectServersCreateRequest) (Server, error) {
	body, code, err := c.request(ctx, "POST", fmt.Sprintf("projects/%d/servers", projectId), data)
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

func (c *Client) ProjectServersAll(ctx context.Context, projectId int) ([]Server, error) {
	resp, err := c.ProjectServers(ctx, projectId)
	if err != nil {
		return nil, err
	}

	servers := resp.Data
	nextPageUrl := resp.Links.Next
	for nextPageUrl != "" {
		body, code, err := c.request(ctx, "GET", nextPageUrl, nil)
		if err != nil {
			return servers, err
		}

		if code != 200 {
			return servers, fmt.Errorf("HTTP %d: %s", code, body)
		}

		var resp ProjectServersResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return servers, fmt.Errorf("failed to decode '%s': %s", body, err)
		}

		servers = append(servers, resp.Data...)
		if nextPageUrl == resp.Links.Next {
			break
		}
		nextPageUrl = resp.Links.Next
	}

	return servers, nil
}

func (c *Client) ProjectServers(ctx context.Context, projectId int) (ProjectServersResponse, error) {
	body, code, err := c.request(ctx, "GET", fmt.Sprintf("projects/%d/servers", projectId), nil)
	if err != nil {
		return ProjectServersResponse{}, err
	}

	if code != 200 {
		return ProjectServersResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ProjectServersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ProjectServersResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}
