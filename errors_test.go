package solus

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPError_Error(t *testing.T) {
	err := newHTTPError(http.MethodDelete, "some/path", http.StatusBadRequest, []byte("foo"))
	require.EqualError(t, err, "HTTP DELETE some/path returns 400 status code: foo")
}
