package solus

import (
	"context"
	"fmt"
)

type SSHKeysService service

type SSHKey struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

type SSHKeyCreateRequest struct {
	Name   string `json:"name"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
}

type SSHKeysResponse struct {
	paginatedResponse

	Data []SSHKey `json:"data"`
}

type sshKeyResponse struct {
	Data SSHKey `json:"data"`
}

func (s *SSHKeysService) List(ctx context.Context, filter *FilterSSHKeys) (SSHKeysResponse, error) {
	resp := SSHKeysResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "ssh_keys", &resp, withFilter(filter.data))
}

func (s *SSHKeysService) Get(ctx context.Context, id int) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("ssh_keys/%d", id), &resp)
}

func (s *SSHKeysService) Create(ctx context.Context, data SSHKeyCreateRequest) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.create(ctx, "ssh_keys", data, &resp)
}

func (s *SSHKeysService) Delete(ctx context.Context, id int) error {
	return s.client.delete(ctx, fmt.Sprintf("ssh_keys/%d", id))
}
