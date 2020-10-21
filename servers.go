package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ServersService service

const (
	ServerStatusProcessing = "processing"
	ServerStatusRunning    = "running"
	ServerStatusStopped    = "stopped"
)

type ServersResponse struct {
	paginatedResponse

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

// Servers return list of server, filter can be nil
func (s *ServersService) List(ctx context.Context, filter *FilterServers) (ServersResponse, error) {
	resp := ServersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}

	body, code, err := s.client.request(ctx, "GET", "servers", withFilter(filter.data))
	if err != nil {
		return ServersResponse{}, err
	}

	if code != 200 {
		return ServersResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return ServersResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}

func (s *ServersService) Get(ctx context.Context, serverId int) (Server, error) {
	body, code, err := s.client.request(ctx, "GET", fmt.Sprintf("servers/%d", serverId))
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

func (s *ServersService) Restart(ctx context.Context, serverId int) (Task, error) {
	body, code, err := s.client.request(ctx, "POST", fmt.Sprintf("servers/%d/restart", serverId))
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

func (s *ServersService) Delete(ctx context.Context, serverId int) (Task, error) {
	body, code, err := s.client.request(ctx, "DELETE", fmt.Sprintf("servers/%d", serverId))
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
