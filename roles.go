package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type RolesService service

type Role struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"is_default"`
	UsersCount int    `json:"users_count"`
}

type RolesResponse struct {
	Data    []Role        `json:"data"`
	Links   ResponseLinks `json:"links"`
	Meta    ResponseMeta  `json:"meta"`
	err     error
	service *RolesService
	opts    requestOpts
}

func (s *RolesService) List(ctx context.Context) (RolesResponse, error) {
	resp := RolesResponse{}

	body, code, err := s.client.request(ctx, "GET", "roles")
	if err != nil {
		return resp, err
	}

	if code != 200 {
		return resp, fmt.Errorf("HTTP %d: %s", code, body)
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}

func (s *RolesService) GetByName(ctx context.Context, name string) (Role, error) {
	roles, err := s.List(ctx)
	if err != nil {
		return Role{}, err
	}

	for _, role := range roles.Data {
		if role.Name == name {
			return role, nil
		}
	}

	return Role{}, fmt.Errorf("failed to get role by name %q: role not found", name)
}
