package solus

import (
	"context"
	"fmt"
	"net/http"
)

// ProjectMember represents a project member.
type ProjectMember struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Email        string `json:"email"`
	Status       string `json:"status"`
	InviteSentAt string `json:"invite_sent_at"`
	IsOwner      bool   `json:"is_owner"`
}

// ProjectMemberCreateRequest represents project member creation request.
type ProjectMemberCreateRequest struct {
	Email string `json:"email"`
}

// ProjectJoin represents response payload when joining a project.
type ProjectJoin struct {
	Email        string `json:"email"`
	IsRegistered bool   `json:"is_registered"`
}

// ProjectTokenPricingTax represents a tax item in project token pricing.
type ProjectTokenPricingTax struct {
	Label string  `json:"label"`
	Rate  float64 `json:"rate"`
}

// ProjectTokenPricing represents project token pricing details.
type ProjectTokenPricing struct {
	UnitCost                  float64                  `json:"unit_cost"`
	CurrencyCode              string                   `json:"currency_code"`
	CurrencyDecimals          int                      `json:"currency_decimals"`
	CurrencyDecimalsSeparator string                   `json:"currency_decimals_separator"`
	CurrencyPrefix            string                   `json:"currency_prefix"`
	CurrencySuffix            string                   `json:"currency_suffix"`
	CurrencyThousandsSep      string                   `json:"currency_thousands_separator"`
	TaxesInclusive            bool                     `json:"taxes_inclusive"`
	Taxes                     []ProjectTokenPricingTax `json:"taxes"`
}

// ProjectSSHKeyCreateRequest represents request payload for creating a project SSH key.
type ProjectSSHKeyCreateRequest struct {
	Name   string `json:"name"`
	Body   string `json:"body"`
	UserID int    `json:"user_id,omitempty"`
}

// Members lists members for the specified project.
func (s *ProjectsService) Members(ctx context.Context, projectID int) ([]ProjectMember, error) {
	var resp struct {
		Data []ProjectMember `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("projects/%d/members", projectID), &resp)
}

// MemberCreate creates a member in the specified project.
func (s *ProjectsService) MemberCreate(
	ctx context.Context,
	projectID int,
	data ProjectMemberCreateRequest,
) (ProjectMember, error) {
	var resp struct {
		Data ProjectMember `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("projects/%d/members", projectID), data, &resp)
}

// MemberDelete deletes a member from the specified project.
func (s *ProjectsService) MemberDelete(ctx context.Context, projectID, memberID int) error {
	return s.client.syncDelete(ctx, fmt.Sprintf("projects/%d/members/%d", projectID, memberID))
}

// MemberResendInvite resends invitation for the specified project member.
func (s *ProjectsService) MemberResendInvite(ctx context.Context, projectID, memberID int) (ProjectMember, error) {
	var resp struct {
		Data ProjectMember `json:"data"`
	}
	return resp.Data, s.client.create(ctx, fmt.Sprintf("projects/%d/members/%d/resend_invite", projectID, memberID), nil, &resp)
}

// Join joins the specified project.
func (s *ProjectsService) Join(ctx context.Context, projectID int) (ProjectJoin, error) {
	path := fmt.Sprintf("projects/%d/join", projectID)
	body, code, err := s.client.request(ctx, http.MethodPost, path)
	if err != nil {
		return ProjectJoin{}, err
	}
	if code != http.StatusOK {
		return ProjectJoin{}, newHTTPError(http.MethodPost, path, code, body)
	}
	var resp struct {
		Data ProjectJoin `json:"data"`
	}
	return resp.Data, unmarshal(body, &resp)
}

// Leave leaves the specified project.
func (s *ProjectsService) Leave(ctx context.Context, projectID int) error {
	path := fmt.Sprintf("projects/%d/leave", projectID)
	body, code, err := s.client.request(ctx, http.MethodPost, path)
	if err != nil {
		return err
	}
	if code != http.StatusOK {
		return newHTTPError(http.MethodPost, path, code, body)
	}
	return nil
}

// Plans lists plans for the specified project.
func (s *ProjectsService) Plans(ctx context.Context, projectID int, filter *FilterPlans) (PlansResponse, error) {
	resp := PlansResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	path := fmt.Sprintf("projects/%d/plans", projectID)
	if filter == nil {
		return resp, s.client.list(ctx, path, &resp)
	}
	return resp, s.client.list(ctx, path, &resp, withFilter(filter.data))
}

// SSHKeys lists SSH keys for the specified project.
func (s *ProjectsService) SSHKeys(ctx context.Context, projectID int) ([]SSHKey, error) {
	var resp struct {
		Data []SSHKey `json:"data"`
	}
	return resp.Data, s.client.get(ctx, fmt.Sprintf("projects/%d/ssh_keys", projectID), &resp)
}

// SSHKeyCreate creates SSH key for the specified project.
func (s *ProjectsService) SSHKeyCreate(
	ctx context.Context,
	projectID int,
	data ProjectSSHKeyCreateRequest,
) (SSHKey, error) {
	var resp sshKeyResponse
	return resp.Data, s.client.create(ctx, fmt.Sprintf("projects/%d/ssh_keys", projectID), data, &resp)
}

// TokenPricing gets token pricing for the specified project.
// If billing is disabled, zero-value ProjectTokenPricing is returned.
func (s *ProjectsService) TokenPricing(ctx context.Context, projectID int) (ProjectTokenPricing, error) {
	path := fmt.Sprintf("projects/%d/token_pricing", projectID)
	body, code, err := s.client.request(ctx, http.MethodGet, path)
	if err != nil {
		return ProjectTokenPricing{}, err
	}
	if code == http.StatusNoContent {
		return ProjectTokenPricing{}, nil
	}
	if code != http.StatusOK {
		return ProjectTokenPricing{}, newHTTPError(http.MethodGet, path, code, body)
	}
	var resp struct {
		Data ProjectTokenPricing `json:"data"`
	}
	return resp.Data, unmarshal(body, &resp)
}
