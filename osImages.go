package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type OsImagesService service

type OsImage struct {
	Id        int              `json:"id"`
	Name      string           `json:"name"`
	Icon      Icon             `json:"icon"`
	Versions  []OsImageVersion `json:"versions,omitempty"`
	IsDefault bool             `json:"is_default,omitempty"`
}

type OsImageVersion struct {
	Id               int    `json:"id"`
	Position         int    `json:"position"`
	Version          string `json:"version"`
	Url              string `json:"url"`
	CloudInitVersion string `json:"cloud_init_version"`
}

type OsImageRequest struct {
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	IsVisible bool   `json:"is_visible"`
}

type OsImageVersionRequest struct {
	Position         int    `json:"position"`
	Version          string `json:"version"`
	Url              string `json:"url"`
	CloudInitVersion string `json:"cloud_init_version"`
}

type OsImageResponse struct {
	Data OsImage `json:"data"`
}

type OsImageVersionResponse struct {
	Data OsImageVersion `json:"data"`
}

type GetOsImageResponse struct {
	Data  []OsImage     `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

func (s *OsImagesService) List(ctx context.Context) ([]OsImage, error) {
	body, code, err := s.client.request(ctx, "GET", "os_images")
	if err != nil {
		return []OsImage{}, err
	}

	if code != 200 {
		return []OsImage{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp GetOsImageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return []OsImage{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *OsImagesService) Create(ctx context.Context, data OsImageRequest) (OsImage, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", "os_images", withBody(opts))
	if err != nil {
		return OsImage{}, err
	}

	if code != 201 {
		return OsImage{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp OsImageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return OsImage{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}

func (s *OsImagesService) OsImageVersionCreate(ctx context.Context, osImageId int, data OsImageVersionRequest) (OsImageVersion, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := s.client.request(ctx, "POST", fmt.Sprintf("os_images/%d/versions", osImageId), withBody(opts))
	if err != nil {
		return OsImageVersion{}, err
	}

	if code != 201 {
		return OsImageVersion{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp OsImageVersionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return OsImageVersion{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
