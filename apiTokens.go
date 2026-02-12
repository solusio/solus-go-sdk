package solus

import (
	"context"
	"fmt"
)

// APITokensService handles all available methods with API tokens.
type APITokensService service

// APIToken represents an API token.
type APIToken struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
}

// AccessToken represents an access token with credentials.
type AccessToken struct {
	Token       APIToken `json:"token"`
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
}

// APITokenCreateRequest represents available properties for creating a new API token.
type APITokenCreateRequest struct {
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

// APITokenPatchRequest represents available properties for patching an existing API token.
type APITokenPatchRequest struct {
	Name string `json:"name"`
}

// APITokensResponse represents paginated list of API tokens.
// This cursor can be used for iterating over all available API tokens.
type APITokensResponse struct {
	paginatedResponse

	Data []APIToken `json:"data"`
}

type apiTokenResponse struct {
	Data APIToken `json:"data"`
}

type accessTokenResponse struct {
	Data AccessToken `json:"data"`
}

// Create creates new API token.
func (s *APITokensService) Create(ctx context.Context, data APITokenCreateRequest) (AccessToken, error) {
	var resp accessTokenResponse
	return resp.Data, s.client.create(ctx, "api_tokens", data, &resp)
}

// List lists API tokens.
func (s *APITokensService) List(ctx context.Context, filter *FilterAPITokens) (APITokensResponse, error) {
	resp := APITokensResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "api_tokens", &resp, withFilter(filter.data))
}

// Get gets specified API token.
func (s *APITokensService) Get(ctx context.Context, id string) (APIToken, error) {
	var resp apiTokenResponse
	return resp.Data, s.client.get(ctx, fmt.Sprintf("api_tokens/%s", id), &resp)
}

// Patch patches specified API token.
func (s *APITokensService) Patch(ctx context.Context, id string, data APITokenPatchRequest) (APIToken, error) {
	var resp apiTokenResponse
	return resp.Data, s.client.patch(ctx, fmt.Sprintf("api_tokens/%s", id), data, &resp)
}

// Delete deletes specified API token.
func (s *APITokensService) Delete(ctx context.Context, id string) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("api_tokens/%s", id))
}
