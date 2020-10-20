package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestLocationsService_List(t *testing.T) {
	expected := LocationsResponse{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/locations", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, url.QueryEscape("filter[search]")+"=test", r.URL.Query().Encode())

		b, err := json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	c := createTestClient(t, s.URL)

	p, err := c.Locations.List(context.Background(), (&FilterLocations{}).ByName("test"))
	require.NoError(t, err)
	p.service = nil
	require.Equal(t, expected, p)
}

func TestLocationsService_Get(t *testing.T) {
	expected := Location{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/locations/1", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)

		b, err := json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.Locations.Get(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}

func TestLocationsService_Create(t *testing.T) {
	expected := Location{}
	data := LocationCreateRequest{
		Name:        "name",
		Description: "description",
		IconId:      null.IntFrom(1),
		IsDefault:   false,
		IsVisible:   true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := LocationCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/locations", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, data, d)

		b, err = json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(201)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.Locations.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}

func TestLocationsService_Update(t *testing.T) {
	expected := Location{}
	data := LocationCreateRequest{
		Name:        "name",
		Description: "description",
		IconId:      null.IntFrom(1),
		IsDefault:   false,
		IsVisible:   true,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := LocationCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/locations/1", r.URL.Path)
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, data, d)

		b, err = json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.Locations.Update(context.Background(), 1, data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}

func TestLocationsService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/locations/1", r.URL.Path)
		require.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(204)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	err = c.Locations.Delete(context.Background(), 1)
	require.NoError(t, err)
}
