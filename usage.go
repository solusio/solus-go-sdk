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
	TerminatedAt string          `json:"terminated_at"`
	Usage        []UsageResource `json:"usage"`
}

// UsageResource represents detailed usage information.
type UsageResource struct {
	Date   string  `json:"date"`
	Tokens float64 `json:"tokens"`
}

// UsageResourceCollection represents collection of usage statistics.
type UsageResourceCollection struct {
	UserID        int                   `json:"user_id"`
	BillingUserID int                   `json:"billing_user_id"`
	Resources     []ServerUsageResource `json:"resources"`
}

// UsageResponse represents response with usage statistics.
type UsageResponse struct {
	Data UsageResourceCollection `json:"data"`
}

// Get retrieves usage statistics.
func (s *UsageService) Get(ctx context.Context, filter *FilterUsage) (UsageResponse, error) {
	var resp UsageResponse
	return resp, s.client.list(ctx, "usage", &resp, withFilter(filter.data))
}
