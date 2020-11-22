package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountService_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/account", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusOK, fakeUser)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Account.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, fakeUser, actual)
}
