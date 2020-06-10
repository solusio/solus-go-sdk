package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

func (c *Client) Users(ctx context.Context) ([]User, error) {
	body, code, err := c.request(ctx, "GET", "users", nil)
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

func (c *Client) UserCreate(ctx context.Context, data UserCreateRequest) (User, error) {
	body, code, err := c.request(ctx, "POST", "users", data)
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

func (c *Client) UserDelete(ctx context.Context, userId int) error {
	body, code, err := c.request(ctx, "DELETE", fmt.Sprintf("users/%d", userId), nil)
	if err != nil {
		return err
	}

	if code != 204 {
		return fmt.Errorf("HTTP %d: %s", code, body)
	}

	return nil
}
