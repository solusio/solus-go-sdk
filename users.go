package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type UsersService service

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	// CreatedAt for date in RFC3339Nano format
	CreatedAt string `json:"created_at"`
}

type UsersResponse struct {
	Data []User `json:"data"`
}

type UserCreateRequest struct {
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Roles    []int  `json:"roles,omitempty"`
}

type UserCreateResponse struct {
	Data User `json:"data"`
}

func (s *UsersService) List(ctx context.Context, filter *FilterUsers) ([]User, error) {
	opts := newRequestOpts()
	opts.params = filterToParams(filter.Get())
	body, code, err := s.client.request(ctx, "GET", "users", withParams(opts))
	if err != nil {
		return []User{}, err
	}

	if code != 200 {
		return []User{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp UsersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return []User{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *UsersService) Create(ctx context.Context, data UserCreateRequest) (User, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "users", withBody(opts))
	if err != nil {
		return User{}, err
	}

	if code != 201 {
		return User{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp UserCreateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return User{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *UsersService) Delete(ctx context.Context, userId int) error {
	body, code, err := s.client.request(ctx, "DELETE", fmt.Sprintf("users/%d", userId))
	if err != nil {
		return err
	}

	if code != 204 {
		return wrapError(code, body)
	}

	return nil
}
