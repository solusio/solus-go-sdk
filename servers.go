package solus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	return resp, s.client.list(ctx, "servers", &resp, withFilter(filter.data))
}

func (s *ServersService) Get(ctx context.Context, serverId int) (Server, error) {
	var resp ServerResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("servers/%d", serverId), &resp)
}

func (s *ServersService) Restart(ctx context.Context, serverId int) (Task, error) {
	body, code, err := s.client.request(ctx, http.MethodPost, fmt.Sprintf("servers/%d/restart", serverId))
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(code, body)
	}

	var resp ServerRestartResponse
	return resp.Data, unmarshal(body, &resp)
}

func (s *ServersService) Delete(ctx context.Context, serverId int) (Task, error) {
	body, code, err := s.client.request(ctx, http.MethodDelete, fmt.Sprintf("servers/%d", serverId))
	if err != nil {
		return Task{}, err
	}

	if code != http.StatusOK {
		return Task{}, newHTTPError(code, body)
	}

	var resp ServerDeleteResponse
	if err := unmarshal(body, &resp); err != nil {
		return Task{}, err
	}

	if resp.Data.Id == 0 {
		return Task{}, errors.New("task doesn't have an id")
	}

	return resp.Data, nil
}
