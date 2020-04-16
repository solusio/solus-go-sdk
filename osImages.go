package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type Icon struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

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

type GetOsImageResponse struct {
	Data  []OsImage     `json:"data"`
	Links ResponseLinks `json:"links"`
	Meta  ResponseMeta  `json:"meta"`
}

func (c *Client) GetOsImages(ctx context.Context) ([]OsImage, error) {
	body, code, err := c.request(ctx, "GET", "os_images", nil)
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
