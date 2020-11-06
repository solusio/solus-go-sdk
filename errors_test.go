package solus

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	err := newHTTPError(http.MethodDelete, "some/path", http.StatusBadRequest, []byte("foo"))
	require.EqualError(t, err, "HTTP DELETE some/path returns 400 status code: foo")
}
