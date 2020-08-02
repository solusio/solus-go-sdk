package solus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	UserStatusActive    = "active"
	UserStatusLocked    = "locked"
	UserStatusSuspended = "suspended"
)

type UsersService service

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	// CreatedAt for date in RFC3339Nano format
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

type UsersResponse struct {
	Data    []User        `json:"data"`
	Links   ResponseLinks `json:"links"`
	Meta    ResponseMeta  `json:"meta"`
	err     error
	service *UsersService
	opts    requestOpts
}

type UserCreateRequest struct {
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	Status     string `json:"status,omitempty"`
	LanguageId int    `json:"language_id,omitempty"`
	Roles      []int  `json:"roles,omitempty"`
}

type UserCreateResponse struct {
	Data User `json:"data"`
}

func (s *UsersService) List(ctx context.Context, filter *FilterUsers) (UsersResponse, error) {
	resp := UsersResponse{}

	opts := newRequestOpts()
	opts.params = filterToParams(filter.Get())

	resp.opts = opts
	resp.service = s

	body, code, err := s.client.request(ctx, "GET", "users", withParams(opts))
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

// Next using for retrieving all data entities
//
// Examples:
//
//  ctx, cancelFunc := context.WithTimeout(context.Background(), 30 * time.Second)
//	defer cancelFunc()
//	usersResponse, err := io.Users.List(ctx, solus.NewFilterUsers())
//	if err != err {
//		return err
//	}
//
//	for usersResponse.Next(ctx) {
//		if err := usersResponse.Err(); err != nil {
//			return err
//		}
//	}
//
//  // all your data now in usersResponse.Data
//	for _, u := range usersResponse.Data {
//
//  }
//
func (ur *UsersResponse) Next(ctx context.Context) bool {
	if ur.Links.Next == "" {
		return false
	}

	nextUrl, err := url.Parse(ur.Links.Next)
	if err != nil {
		ur.err = err
		return false
	}

	for k, v := range nextUrl.Query() {
		ur.opts.params[k] = v
	}

	body, code, err := ur.service.client.request(ctx, "GET", "users", withParams(ur.opts))
	if err != nil {
		ur.err = err
		return false
	}

	if code != 200 {
		ur.err = fmt.Errorf("HTTP %d: %s", code, body)
		return false
	}

	var resp UsersResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		ur.err = fmt.Errorf("failed to decode '%s': %s", body, err)
		return false
	}

	ur.Data = append(ur.Data, resp.Data...)
	ur.Meta = resp.Meta
	ur.Links = resp.Links

	return true
}

func (ur *UsersResponse) Err() error {
	return ur.err
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
