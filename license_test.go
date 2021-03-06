package solus

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLicenseService_Activate(t *testing.T) {
	expected := fakeLicense
	data := LicenseActivateRequest{
		ActivationCode: "activation code",
	}

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
}
