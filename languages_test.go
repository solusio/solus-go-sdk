package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeLanguage = Language{
	ID:         1,
	Name:       "English",
	Locale:     "en_US",
	Country:    "US",
	Icon:       fakeIcon,
	IsDefault:  true,
	IsVisible:  true,
	UsersCount: 10,
}

func TestLanguagesService_List(t *testing.T) {
	expected := LanguagesResponse{
		Data: []Language{fakeLanguage},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/languages", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"English"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterLanguages{}).ByName("English")

	actual, err := createTestClient(t, s.URL).Languages.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}
