package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestLicenseService_Activate(t *testing.T) {
	data := LicenseActivateRequest{
		ActivationCode: "activation code",
	}

	t.Run("positive", func(t *testing.T) {
		expected := fakeLicense

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/license/activate", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)
			assertRequestBody(t, r, data)

			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		actual, err := createTestClient(t, s.URL).License.Activate(context.Background(), data)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			asserter, addr := startBrokenTestServer(t)
			_, err := createTestClient(t, addr).License.Activate(context.Background(), data)
			asserter(t, http.MethodPost, "/license/activate", err)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/license/activate", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)
				assertRequestBody(t, r, data)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).License.Activate(context.Background(), data)
			assert.EqualError(t, err, "HTTP POST license/activate returns 400 status code")
		})
	})
}
