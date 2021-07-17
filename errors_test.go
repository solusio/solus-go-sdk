package solus

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPError_Error(t *testing.T) {
	t.Run("a JSON body", func(t *testing.T) {
		err := newHTTPError(
			http.MethodDelete,
			"some/path",
			http.StatusBadRequest,
			[]byte(`{"message": "foo"}`),
		)
		require.EqualError(t, err, "HTTP DELETE some/path returns 400 status code: foo")
	})

	t.Run("not a JSON body", func(t *testing.T) {
		err := newHTTPError(http.MethodDelete, "some/path", http.StatusBadRequest, []byte("foo"))
		require.EqualError(t, err, "HTTP DELETE some/path returns 400 status code: foo")
	})
}

func TestIsNotFound(t *testing.T) {
	testCases := map[string]struct {
		err      error
		expected bool
	}{
		"not http err": {
			errors.New("fake error"),
			false,
		},
		"http error, not 404": {
			newHTTPError(http.MethodPut, "/foo", http.StatusBadRequest, nil),
			false,
		},
		"404 http error": {
			newHTTPError(http.MethodPut, "/foo", http.StatusNotFound, nil),
			true,
		},
	}

	for name, tt := range testCases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsNotFound(tt.err))
		})
	}
}
