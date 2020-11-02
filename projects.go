package solus

import (
	"context"
)

type ProjectsService service

type ProjectsResponse struct {
	paginatedResponse

	Data []Project `json:"data"`
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
	resp := ProjectsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "projects", &resp)
}
