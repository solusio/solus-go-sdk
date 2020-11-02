package solus

import (
	"context"
	"net/http"
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
	body, code, err := s.client.request(ctx, http.MethodPost, "license/activate", withBody(data))
	if err != nil {
		return License{}, err
	}

	if code != http.StatusOK {
		return License{}, newHTTPError(code, body)
	}

	var resp LicenseActivateResponse
	return resp.Data, unmarshal(body, &resp)
}
