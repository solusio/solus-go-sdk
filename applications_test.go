package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
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
		Url:              "url",
		IconID:           1,
		CloudInitVersion: "cloud init version",
		UserDataTemplate: "user data template",
		JsonSchema:       "json schema",
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
