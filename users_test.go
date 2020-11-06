package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestUsersService_List(t *testing.T) {
	expected := UsersResponse{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[status]": []string{"status"},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterUsers{}).ByStatus("status")

	actual, err := createTestClient(t, s.URL).Users.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestUsersService_Create(t *testing.T) {
	data := UserCreateRequest{
		Password:   "password",
		Email:      "email",
		Status:     "status",
		LanguageID: 1,
		Roles:      []int{2, 3},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeUser)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Users.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, fakeUser, actual)
}

func TestUsersService_Update(t *testing.T) {
	data := UserUpdateRequest{
		Password:   "password",
		Status:     "status",
		LanguageID: 1,
		Roles:      []int{2, 3},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/10", r.URL.Path)
		assert.Equal(t, http.MethodPut, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusOK, fakeUser)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Users.Update(context.Background(), 10, data)
	require.NoError(t, err)
	require.Equal(t, fakeUser, actual)
}

func TestUsersService_Delete(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users/10", r.URL.Path)
		assert.Equal(t, http.MethodDelete, r.Method)

		w.WriteHeader(http.StatusNoContent)
	})
	defer s.Close()

	err := createTestClient(t, s.URL).Users.Delete(context.Background(), 10)
	require.NoError(t, err)
}
