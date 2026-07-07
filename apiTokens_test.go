package solus

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeAPIToken = APIToken{
	ID:        "token-1",
	Name:      "primary token",
	User:      fakeUser,
	CreatedAt: "2026-02-11T15:46:30Z",
}

var fakeAccessToken = AccessToken{
	Token:       fakeAPIToken,
	AccessToken: "secret-token-value",
	TokenType:   "Bearer",
}

func TestAPITokensService_List(t *testing.T) {
	expected := APITokensResponse{
		Data: []APIToken{fakeAPIToken},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api_tokens", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[search]":  []string{"primary"},
			"filter[user_id]": []string{"1"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterAPITokens{}).
		ByName("primary").
		ByUserID(1)

	actual, err := createTestClient(t, s.URL).APITokens.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestAPITokensService_Create(t *testing.T) {
	req := APITokenCreateRequest{
		Name:   "primary token",
		UserID: 1,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api_tokens", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusCreated, fakeAccessToken)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).APITokens.Create(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, fakeAccessToken, actual)
}

func TestAPITokensService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api_tokens/token-1", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeAPIToken)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).APITokens.Get(context.Background(), "token-1")
	require.NoError(t, err)
	require.Equal(t, fakeAPIToken, actual)
}

func TestAPITokensService_Patch(t *testing.T) {
	req := APITokenPatchRequest{Name: "renamed token"}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api_tokens/token-1", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)
		assertRequestBody(t, r, req)

		writeResponse(t, w, http.StatusOK, fakeAPIToken)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).APITokens.Patch(context.Background(), "token-1", req)
	require.NoError(t, err)
	require.Equal(t, fakeAPIToken, actual)
}

func TestAPITokensService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api_tokens/token-1", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).APITokens.Delete(context.Background(), "token-1")
	require.NoError(t, err)
}
