package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ProjectsService service

type ProjectsResponse struct {
	Data  []Project     `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

type Project struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     int      `json:"members"`
	IsOwner     bool     `json:"is_owner"`
	IsDefault   bool     `json:"is_default"`
	Owner       User     `json:"owner"`
	Servers     []Server `json:"servers"`
}

func (s *ProjectsService) List(ctx context.Context) (ProjectsResponse, error) {
	body, code, err := s.client.request(ctx, "GET", "projects")
	if err != nil {
		return ProjectsResponse{}, err
	}

	if code != 200 {
		return ProjectsResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ProjectsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return ProjectsResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp, nil
}
