package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRolesService_List(t *testing.T) {
	expected := RolesResponse{
		Data: []Role{
			fakeRole,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/roles", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Roles.List(context.Background())
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}

func TestRolesService_GetByName(t *testing.T) {
	resp := RolesResponse{
		Data: []Role{
			{Name: "foo"},
			{Name: "bar"},
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/roles", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeJSON(t, w, http.StatusOK, resp)
	})
	defer s.Close()

	t.Run("positive", func(t *testing.T) {
		actual, err := createTestClient(t, s.URL).Roles.GetByName(context.Background(), "foo")
		require.NoError(t, err)
		require.Equal(t, Role{Name: "foo"}, actual)
	})

	t.Run("positive", func(t *testing.T) {
		_, err := createTestClient(t, s.URL).Roles.GetByName(context.Background(), "fizz")
		require.EqualError(t, err, `failed to get role by name "fizz": role not found`)
	})
}
