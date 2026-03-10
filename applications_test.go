package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationsService_List(t *testing.T) {
	expected := ApplicationsResponse{
		Data: []Application{
			fakeApplication,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestApplicationsService_Create(t *testing.T) {
	data := ApplicationCreateRequest{
		Name:             "name",
		URL:              "url",
		IconID:           1,
		CloudInitVersion: "cloud init version",
		UserDataTemplate: "user data template",
		JSONSchema:       "json schema",
		IsVisible:        true,
		LoginLink: LoginLink{
			Type:    LoginLinkTypeURL,
			Content: "login link content",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeApplication)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeApplication, actual)
}

func TestApplicationsService_Update(t *testing.T) {
	data := ApplicationCreateRequest{
		Name:             "updated-name",
		URL:              "https://example.com/app",
		IconID:           10,
		CloudInitVersion: "v3",
		UserDataTemplate: "template",
		JSONSchema:       "{\"type\":\"object\"}",
		IsVisible:        false,
		LoginLink: LoginLink{
			Type:    LoginLinkTypeURL,
			Content: "https://example.com/login",
		},
		AvailablePlans: []int{1, 2},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/applications/42", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeApplication)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Applications.Update(context.Background(), 42, data)
	require.NoError(t, err)
	require.Equal(t, fakeApplication, actual)
}
