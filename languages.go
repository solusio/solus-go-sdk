package solus

import (
	"context"
)

// LanguagesService handles all available methods with languages.
type LanguagesService service

// Language represents a language.
type Language struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Locale     string `json:"locale"`
	Country    string `json:"country"`
	Icon       Icon   `json:"icon"`
	IsDefault  bool   `json:"is_default"`
	IsVisible  bool   `json:"is_visible"`
	UsersCount int    `json:"users_count"`
}

// LanguagesResponse represents paginated list of languages.
// This cursor can be used for iterating over all available languages.
type LanguagesResponse struct {
	paginatedResponse

	Data []Language `json:"data"`
}

// List lists languages.
func (s *LanguagesService) List(ctx context.Context, filter *FilterLanguages) (LanguagesResponse, error) {
	resp := LanguagesResponse{
		paginatedResponse: paginatedResponse{
			service: (*service)(s),
		},
	}
	return resp, s.client.list(ctx, "languages", &resp, withFilter(filter.data))
}
