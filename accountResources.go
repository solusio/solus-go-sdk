package solus

import (
	"context"
)

// Notification represents account notification.
type Notification struct {
	ID   int                    `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// AccountSSHKeyCreateRequest represents request payload for creating account SSH key.
type AccountSSHKeyCreateRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

// AccountTokenCreateRequest represents request payload for creating account API token.
type AccountTokenCreateRequest struct {
	Name string `json:"name"`
}

// Notifications lists current account notifications.
func (s *AccountService) Notifications(ctx context.Context) ([]Notification, error) {
	var resp struct {
		Data []Notification `json:"data"`
	}
	return resp.Data, s.client.get(ctx, "account/notifications", &resp)
}

// NotificationsDelete deletes all current account notifications.
func (s *AccountService) NotificationsDelete(ctx context.Context) error {
	return s.client.syncDelete(ctx, "account/notifications")
}

// SSHKeys lists current account SSH keys.
func (s *AccountService) SSHKeys(ctx context.Context) (SSHKeysResponse, error) {
	resp := SSHKeysResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "account/ssh_keys", &resp)
}

// SSHKeyCreate creates a new current account SSH key.
func (s *AccountService) SSHKeyCreate(ctx context.Context, data AccountSSHKeyCreateRequest) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.create(ctx, "account/ssh_keys", data, &resp)
}

// Tokens lists current account API tokens.
func (s *AccountService) Tokens(ctx context.Context) (APITokensResponse, error) {
	resp := APITokensResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "account/tokens", &resp)
}

// TokenCreate creates a new current account API token.
func (s *AccountService) TokenCreate(ctx context.Context, data AccountTokenCreateRequest) (AccessToken, error) {
	var resp accessTokenResponse
	return resp.Data, s.client.create(ctx, "account/tokens", data, &resp)
}
