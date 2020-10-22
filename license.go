package solus

import (
	"context"
	"encoding/json"
	"fmt"
)

type LicenseService service

type License struct {
	CpuCores       int    `json:"cpu_cores"`
	CpuCoresInUse  int    `json:"cpu_cores_in_use"`
	IsActive       bool   `json:"is_active"`
	Key            string `json:"key"`
	KeyType        string `json:"key_type"`
	Product        string `json:"product"`
	ExpirationDate string `json:"expiration_date"`
	UpdateDate     string `json:"update_date"`
}

type LicenseActivateRequest struct {
	ActivationCode string `json:"activation_code"`
}

type LicenseActivateResponse struct {
	Data License `json:"data"`
}

func (s *LicenseService) Activate(ctx context.Context, data LicenseActivateRequest) (License, error) {
	body, code, err := s.client.request(ctx, "POST", "license/activate", withBody(data))
	if err != nil {
		return License{}, err
	}

	if code != 200 {
		return License{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp LicenseActivateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return License{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
