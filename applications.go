package solus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// ApplicationsService handles all available methods with applications.
type ApplicationsService service

// Application represents an application.
type Application struct {
	ID                  int                             `json:"id"`
	Name                string                          `json:"name"`
	Icon                Icon                            `json:"icon"`
	URL                 string                          `json:"url"`
	CloudInitVersion    string                          `json:"cloud_init_version"`
	UserData            string                          `json:"user_data_template"`
	LoginLink           LoginLink                       `json:"login_link"`
	JSONSchema          string                          `json:"json_schema"`
	IsDefault           bool                            `json:"is_default"`
	IsVisible           bool                            `json:"is_visible"`
	IsBuiltin           bool                            `json:"is_buildin"`
	AvailablePlans      []ShortPlan                     `json:"available_plans"`
	VSDiskDriver        *string                         `json:"vs_disk_driver,omitempty"`
	VSFirmware          *string                         `json:"vs_firmware,omitempty"`
	IsDeprecated        bool                            `json:"is_deprecated,omitempty"`
	ExternalIntegration *ApplicationExternalIntegration `json:"external_integration,omitempty"`
	HasVSDiskDriver     bool                            `json:"-"`
	HasVSFirmware       bool                            `json:"-"`
}

// ApplicationExternalIntegration identifies an application's external integration.
type ApplicationExternalIntegration struct {
	ID int `json:"id"`
}

// UnmarshalJSON preserves the distinction between missing and nullable settings.
func (a *Application) UnmarshalJSON(data []byte) error {
	type plain Application
	var decoded plain
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	_, decoded.HasVSDiskDriver = fields["vs_disk_driver"]
	_, decoded.HasVSFirmware = fields["vs_firmware"]
	*a = Application(decoded)
	return nil
}

// ShortApplication represents only ID and name of application.
type ShortApplication struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// LoginLinkType a type of login link to the application.
type LoginLinkType string

const (
	// LoginLinkTypeNone indicates application without login link.
	LoginLinkTypeNone LoginLinkType = "none"

	// LoginLinkTypeURL indicates application with URL pattern login link.
	LoginLinkTypeURL LoginLinkType = "url"

	// LoginLinkTypeJSCode indicates application with custom JS code auth code.
	LoginLinkTypeJSCode LoginLinkType = "js_code"

	// LoginLinkTypeJSCodeCanonical is the API spelling for JavaScript login links.
	LoginLinkTypeJSCodeCanonical LoginLinkType = "js-code"

	// LoginLinkTypeInfo indicates application with custom information in popup.
	LoginLinkTypeInfo LoginLinkType = "info"
)

// LoginLink represents an application login link.
type LoginLink struct {
	Type    LoginLinkType `json:"type"`
	Content string        `json:"content"`
}

// ApplicationCreateRequest represents available properties for creating a new
// application.
type ApplicationCreateRequest struct {
	Name             string    `json:"name"`
	URL              string    `json:"url"`
	IconID           int       `json:"icon_id"`
	CloudInitVersion string    `json:"cloud_init_version"`
	UserDataTemplate string    `json:"user_data_template"`
	JSONSchema       string    `json:"json_schema"`
	IsVisible        bool      `json:"is_visible"`
	LoginLink        LoginLink `json:"login_link"`
	AvailablePlans   []int     `json:"available_plans,omitempty"`
}

// ApplicationsResponse represents paginated list of applications.
// This cursor can be used for iterating over all available applications.
type ApplicationsResponse struct {
	paginatedResponse

	Data []Application `json:"data"`
}

// Create creates new application.
func (s *ApplicationsService) Create(ctx context.Context, data ApplicationCreateRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	return resp.Data, s.client.create(ctx, "applications", data, &resp)
}

// Update updates an existing application.
func (s *ApplicationsService) Update(ctx context.Context, id int, data ApplicationCreateRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	return resp.Data, s.client.update(ctx, fmt.Sprintf("applications/%d", id), data, &resp)
}

// List lists all applications.
func (s *ApplicationsService) List(ctx context.Context) (ApplicationsResponse, error) {
	resp := ApplicationsResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "applications", &resp)
}

// ApplicationWriteRequest includes explicit empty and nullable settings for full writes.
type ApplicationWriteRequest struct {
	ApplicationCreateRequest
	AvailablePlans        []int   `json:"available_plans"`
	VSDiskDriver          *string `json:"vs_disk_driver"`
	VSFirmware            *string `json:"vs_firmware"`
	ExternalIntegrationID *int    `json:"external_integration_id"`
}

// Get retrieves an application by ID.
func (s *ApplicationsService) Get(ctx context.Context, id int) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	err := s.client.get(ctx, fmt.Sprintf("applications/%d", id), &resp)
	return resp.Data, applicationOperationError(err)
}

// Delete removes an application by ID.
func (s *ApplicationsService) Delete(ctx context.Context, id int) error {
	return applicationOperationError(s.client.syncDelete(ctx, fmt.Sprintf("applications/%d", id)))
}

// CreateFull creates an application with explicit full-write settings.
func (s *ApplicationsService) CreateFull(ctx context.Context, data ApplicationWriteRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	err := s.client.create(ctx, "applications", data, &resp)
	return resp.Data, applicationOperationError(err)
}

// UpdateFull updates all writable application settings.
func (s *ApplicationsService) UpdateFull(ctx context.Context, id int, data ApplicationWriteRequest) (Application, error) {
	var resp struct {
		Data Application `json:"data"`
	}
	err := s.client.update(ctx, fmt.Sprintf("applications/%d", id), data, &resp)
	return resp.Data, applicationOperationError(err)
}

func applicationOperationError(err error) error {
	if err == nil {
		return nil
	}
	var syntax *json.SyntaxError
	if errors.As(err, &syntax) {
		return fmt.Errorf("decode application response: %w", syntax)
	}
	var field *json.UnmarshalTypeError
	if errors.As(err, &field) {
		sanitized := *field
		sanitized.Value = "invalid value"
		return fmt.Errorf("decode application response: %w", &sanitized)
	}
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		httpErr.Message = "application API request failed"
		for key := range httpErr.Errors {
			httpErr.Errors[key] = []string{"invalid value"}
		}
		return httpErr
	}
	return err
}
