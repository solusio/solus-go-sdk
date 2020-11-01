package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"net/http"
	"net/url"
	"testing"
)

func TestLocationsService_Create(t *testing.T) {
	expected := fakeLocation
	data := LocationCreateRequest{
		Name:        "name",
		Description: "description",
		IconId:      null.IntFrom(1),
		IsDefault:   false,
		IsVisible:   true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/locations", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Locations.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestLocationsService_List(t *testing.T) {
	expected := LocationsResponse{
		Data: []Location{
			fakeLocation,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/locations", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]": []string{"name"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterLocations{}).ByName("name")

	actual, err := createTestClient(t, s.URL).Locations.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestLocationsService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/locations/10", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeLocation)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Locations.Get(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, fakeLocation, actual)
}

func TestLocationsService_Update(t *testing.T) {
	data := LocationCreateRequest{
		Name:        "name",
		Description: "description",
		IconId:      null.IntFrom(1),
		IsDefault:   false,
		IsVisible:   true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/locations/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeLocation)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Locations.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeLocation, actual)
}

func TestLocationsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/locations/10", r.URL.Path)
		require.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Locations.Delete(context.Background(), 10)
	require.NoError(t, err)
}
