package solus

import (
	"context"
)

// UsageService handles all available methods with usage statistics.
type UsageService service

// ServerUsageResource represents server usage information.
type ServerUsageResource struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	TotalTokens  float64         `json:"total_tokens"`
	CreatedAt    string          `json:"created_at"`
	TerminatedAt *string         `json:"terminated_at"`
	Usage        []UsageResource `json:"usage"`
}

// UsageResource represents detailed usage information.
type UsageResource struct {
	Type               string  `json:"type"`
	Description        string  `json:"description"`
	Quantity           float64 `json:"quantity"`
	QuantityDescriptor string  `json:"quantity_descriptor"`
	Tokens             int     `json:"tokens"`
	TokensPerHour      int     `json:"tokens_per_hour"`
	TokensPerMonth     int     `json:"tokens_per_month"`
}

// UsageResourceCollection represents usage statistics grouped by user.
type UsageResourceCollection struct {
	UserID        int                   `json:"user_id"`
	BillingUserID *int                  `json:"billing_user_id"`
	Resources     []ServerUsageResource `json:"resources"`
}

// UsageResponse represents paginated usage statistics response.
type UsageResponse struct {
	paginatedResponse

	Data []UsageResourceCollection `json:"data"`
}

// Get retrieves usage statistics.
func (s *UsageService) Get(ctx context.Context, filter *FilterUsage) (UsageResponse, error) {
	resp := UsageResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "usage", &resp, withFilter(filter.data))
}
