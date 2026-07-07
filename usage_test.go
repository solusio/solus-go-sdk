package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsageService_Get(t *testing.T) {
	expected := UsageResponse{
		paginatedResponse: paginatedResponse{
			Links: ResponseLinks{
				First: "https://example.test/usage?page=1",
				Last:  "https://example.test/usage?page=1",
			},
			Meta: ResponseMeta{
				CurrentPage: 1,
				LastPage:    1,
				PerPage:     25,
				Total:       1,
			},
		},
		Data: []UsageResourceCollection{
			{
				UserID:        7,
				BillingUserID: nil,
				Resources: []ServerUsageResource{
					{
						ID:           52494,
						Name:         "virtual server",
						TotalTokens:  0,
						CreatedAt:    "2026-01-01 00:00:14",
						TerminatedAt: nil,
						Usage: []UsageResource{
							{
								Type:               "up-time",
								Description:        "whmcs.source (1 vCPU, 1 GiB RAM, 10 GiB disk)",
								Quantity:           744,
								QuantityDescriptor: "hours",
								Tokens:             0,
								TokensPerHour:      0,
								TokensPerMonth:     0,
							},
						},
					},
				},
			},
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/usage", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[period]":          []string{"2026-02"},
			"filter[user_id]":         []string{"7"},
			"filter[billing_user_id]": []string{"9"},
			"timezone":                []string{"UTC"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	filter := (&FilterUsage{}).
		ByPeriod("2026-02").
		ByUserID(7).
		ByBillingUserID(9).
		WithTimezone("UTC")

	actual, err := createTestClient(t, s.URL).Usage.Get(context.Background(), filter)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}
