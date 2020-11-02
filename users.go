package solus

import (
	"context"
	"fmt"
)

const (
	UserStatusActive    = "active"
	UserStatusLocked    = "locked"
	UserStatusSuspended = "suspended"
)

type UsersService service

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// CreatedAt for date in RFC3339Nano format
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Roles     []Role `json:"roles"`
}

type UsersResponse struct {
	paginatedResponse

	Data []User `json:"data"`
}

type UserCreateRequest struct {
	Password   string `json:"password,omitempty"`
	Email      string `json:"email,omitempty"`
	Status     string `json:"status,omitempty"`
	LanguageId int    `json:"language_id,omitempty"`
	Roles      []int  `json:"roles,omitempty"`
}

type UserUpdateRequest struct {
	Password   string `json:"password,omitempty"`
	Status     string `json:"status,omitempty"`
	LanguageId int    `json:"language_id,omitempty"`
	Roles      []int  `json:"roles,omitempty"`
}

type UserCreateResponse struct {
	Data User `json:"data"`
}

func (s *UsersService) List(ctx context.Context, filter *FilterUsers) (UsersResponse, error) {
	resp := UsersResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "users", &resp, withFilter(filter.data))
}

func (s *UsersService) Create(ctx context.Context, data UserCreateRequest) (User, error) {
	var resp UserCreateResponse
	return resp.Data, s.client.create(ctx, "users", data, &resp)
}

func (s *UsersService) Update(ctx context.Context, userId int, data UserUpdateRequest) (User, error) {
	var resp UserCreateResponse
	return resp.Data, s.client.update(ctx, fmt.Sprintf("users/%d", userId), data, &resp)
}

func (s *UsersService) Delete(ctx context.Context, userId int) error {
	return s.client.delete(ctx, fmt.Sprintf("users/%d", userId))
}
