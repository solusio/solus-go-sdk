package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIconsService_List(t *testing.T) {
	expected := IconsResponse{
		Data: []Icon{
			fakeIcon,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/icons", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
			"filter[type]":   []string{string(IconTypeApplication)},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterIcons{}).
		ByName("name").
		ByType(IconTypeApplication)

	actual, err := createTestClient(t, s.URL).Icons.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestIconsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/icons/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeIcon)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Icons.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeIcon, actual)
}
