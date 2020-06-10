package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	ServerStatusProcessing = "processing"
	ServerStatusRunning    = "running"
	ServerStatusStopped    = "stopped"
)

type ServersResponse struct {
	Data []Server `json:"data"`
}

type ServerResponse struct {
	Data Server `json:"data"`
}

type ServerRestartResponse struct {
	Data Task `json:"data"`
}

type ServerDeleteResponse struct {
	Data Task `json:"data"`
}

func (c *Client) Servers(ctx context.Context) ([]Server, error) {
	body, code, err := c.request(ctx, "GET", "servers", nil)
	if err != nil {
		return []Server{}, err
	}

	if code != 200 {
		return []Server{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ServersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return []Server{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) Server(ctx context.Context, serverId int) (Server, error) {
	body, code, err := c.request(ctx, "GET", fmt.Sprintf("servers/%d", serverId), nil)
	if err != nil {
		return Server{}, err
	}

	if code != 200 {
		return Server{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ServerResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Server{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) ServerRestart(ctx context.Context, serverId int) (Task, error) {
	body, code, err := c.request(ctx, "POST", fmt.Sprintf("servers/%d/restart", serverId), nil)
	if err != nil {
		return Task{}, err
	}

	if code != 200 {
		return Task{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ServerRestartResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Task{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (c *Client) ServerDelete(ctx context.Context, serverId int) (Task, error) {
	body, code, err := c.request(ctx, "DELETE", fmt.Sprintf("servers/%d", serverId), nil)
	if err != nil {
		return Task{}, err
	}

	if code != 200 {
		return Task{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ServerDeleteResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Task{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	if resp.Data.Id == 0 {
		return Task{}, fmt.Errorf("failed to decode '%s': to task", body)
	}

	return resp.Data, nil
}
