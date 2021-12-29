package solus

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPError_Error(t *testing.T) {
	cc := map[string]struct {
		body            string
		expectedMessage string
		expectedErrors  map[string][]string
	}{
		"empty body": {
			body:            "",
			expectedMessage: "HTTP DELETE some/path returns 400 status code",
		},

		"not a JSON body": {
			body:            "foo",
			expectedMessage: "HTTP DELETE some/path returns 400 status code: foo",
		},

		"empty message": {
			body: `{
	"message": ""
}`,
			expectedMessage: "HTTP DELETE some/path returns 400 status code",
		},

		"with message only": {
			body: `{
	"message": "foo"
}`,
			expectedMessage: "HTTP DELETE some/path returns 400 status code: foo",
		},

		"with errors only": {
			body: `{
	"errors": {
		"foo": [
			"fizz",
			"buzz"
		],
		"bar": [
			"foo"
		]
	}
}`,
			expectedMessage: `HTTP DELETE some/path returns 400 status code with errors`,
			expectedErrors: map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			},
		},

		"with message and errors": {
			body: `{
	"message": "foo",
	"errors": {
		"foo": [
			"fizz",
			"buzz"
		],
		"bar": [
			"foo"
		]
	}
}`,
			expectedMessage: `HTTP DELETE some/path returns 400 status code with errors: foo`,
			expectedErrors: map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			},
		},
	}

	for n, c := range cc {
		t.Run(n, func(t *testing.T) {
			err := newHTTPError(
				http.MethodDelete,
				"some/path",
				http.StatusBadRequest,
				[]byte(c.body),
			)
			require.EqualError(t, err, c.expectedMessage)

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, c.expectedErrors, e.Errors)
		})
	}
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
