package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type ApplicationsService service

type LoginLinkType string

const (
	LoginLinkTypeNone   LoginLinkType = "none"
	LoginLinkTypeURL    LoginLinkType = "url"
	LoginLinkTypeJSCode LoginLinkType = "js_code"
	LoginLinkTypeInfo   LoginLinkType = "info"
)

type LoginLink struct {
	Type    LoginLinkType `json:"type"`
	Content string        `json:"content"`
}

type Application struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Icon             Icon      `json:"icon"`
	Url              string    `json:"url"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserData         string    `json:"user_data_template"`
	LoginLink        LoginLink `json:"login_link"`
	JsonSchema       string    `json:"json_schema"`
	IsDefault        bool      `json:"is_default"`
	IsVisible        bool      `json:"is_visible"`
	IsBuiltin        bool      `json:"is_buildin"`
}

type ApplicationRequest struct {
	Name             string    `json:"name"`
	Url              string    `json:"url"`
	IconId           int       `json:"icon_id"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserDataTemplate string    `json:"user_data_template"`
	JsonSchema       string    `json:"json_schema"`
	IsVisible        bool      `json:"is_visible"`
	LoginLink        LoginLink `json:"login_link"`
}

type ApplicationResponse struct {
	Data Application `json:"data"`
}

type ApplicationsPaginatedResponse struct {
	Data  []Application `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

func (s *ApplicationsService) List(ctx context.Context) ([]Application, error) {
	body, code, err := s.client.request(ctx, "GET", "applications")
	if err != nil {
		return []Application{}, err
	}

	if code != 200 {
		return []Application{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ApplicationsPaginatedResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return []Application{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *ApplicationsService) Create(ctx context.Context, data ApplicationRequest) (Application, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "applications", withBody(opts))
	if err != nil {
		return Application{}, err
	}

	if code != 201 {
		return Application{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp ApplicationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return Application{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
