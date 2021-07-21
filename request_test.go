package solus

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func Test_withFilter(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		o := requestOpts{}

		withFilter(nil)(&o)

		assert.Equal(t, map[string][]string(nil), o.params)
	})

	t.Run("without params", func(t *testing.T) {
		o := requestOpts{}

		withFilter(map[string]string{
			"foo":  "bar",
			"fizz": "buzz",
		})(&o)

		assert.Equal(t, map[string][]string{
			"foo":  {"bar"},
			"fizz": {"buzz"},
		}, o.params)
	})

	t.Run("with params", func(t *testing.T) {
		o := requestOpts{
			params: map[string][]string{
				"foo": {"new"},
				"100": {"500"},
			},
		}

		withFilter(map[string]string{
			"foo":  "bar",
			"fizz": "buzz",
		})(&o)

		assert.Equal(t, map[string][]string{
			"foo":  {"new", "bar"},
			"fizz": {"buzz"},
			"100":  {"500"},
		}, o.params)
	})
}

func TestClient_buildRequest(t *testing.T) {
	const (
		rawBaseURL     = "http://foo"
		additionalPath = "fizz/buzz"
		userAgent      = "user agent"
	)

	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)

	cl := &Client{
		UserAgent: userAgent,
		BaseURL:   baseURL,
		Logger:    NullLogger{},
		Headers: http.Header{
			"foo": {"bar"},
		},
	}

	t.Run("positive", func(t *testing.T) {
		t.Run("without options", func(t *testing.T) {
			expectedURL, err := url.Parse(fmt.Sprintf("%s/%s", rawBaseURL, additionalPath))
			require.NoError(t, err)

			r, err := cl.buildRequest(context.Background(), http.MethodPost, additionalPath)

			require.NoError(t, err)
			assert.Equal(t, expectedURL, r.URL)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, map[string][]string{
				"Foo":        {"bar"},
				"User-Agent": {userAgent},
			}, map[string][]string(r.Header))
			assert.Nil(t, r.Body)
		})

		t.Run("with options", func(t *testing.T) {
			expectedURL, err := url.Parse(fmt.Sprintf(
				"%s/%s?foo=bar&fizz=buzz",
				rawBaseURL,
				additionalPath,
			))
			require.NoError(t, err)

			body := map[string]interface{}{
				"a": "b",
			}

			r, err := cl.buildRequest(
				context.Background(),
				http.MethodPost,
				additionalPath,
				withFilter(map[string]string{
					"foo":  "bar",
					"fizz": "buzz",
				}),
				withBody(body),
			)

			require.NoError(t, err)

			// Because we didn't have any guaranties about order of query parameters.
			expectedQuery := expectedURL.Query()
			actualQuery := r.URL.Query()
			expectedURL.RawQuery = ""
			r.URL.RawQuery = ""

			assert.Equal(t, expectedURL, r.URL)
			assert.Equal(t, expectedQuery, actualQuery)
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, map[string][]string{
				"Foo":        {"bar"},
				"User-Agent": {userAgent},
			}, map[string][]string(r.Header))

			assertRequestBody(t, r, body)
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to create request", func(t *testing.T) {
			_, err := cl.buildRequest(
				context.Background(),
				"foo bar",
				"",
			)

			assert.EqualError(t, err, `net/http: invalid method "foo bar"`)
		})
	})
}

func TestClient_buildURL(t *testing.T) {
	const (
		rawBaseURL     = "http://foo"
		additionalPath = "fizz/buzz"
	)
	expectedURL := fmt.Sprintf("%s/%s", rawBaseURL, additionalPath)

	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)

	cl := &Client{
		BaseURL: baseURL,
	}

	t.Run("without parameters", func(t *testing.T) {
		u, err := cl.buildURL(additionalPath, requestOpts{})
		require.NoError(t, err)
		assert.Equal(t, expectedURL, u)
	})

	t.Run("with parameters", func(t *testing.T) {
		expectedParams := map[string][]string{
			"foo":  {"bar"},
			"fizz": {"buzz", "42"},
		}

		u, err := cl.buildURL(additionalPath, requestOpts{
			params: expectedParams,
		})
		require.NoError(t, err)

		actual, err := url.Parse(u)
		require.NoError(t, err)

		assert.Equal(t, expectedParams, map[string][]string(actual.Query()))
	})
}

func Test_checkForRetry(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		ok, err := checkForRetry(&http.Response{StatusCode: 200}, nil)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil response", func(t *testing.T) {
			assert.Panics(t, func() {
				_, _ = checkForRetry(nil, nil)
			})
		})

		t.Run("have an error", func(t *testing.T) {
			_, err := checkForRetry(nil, errors.New("fake error"))
			assert.EqualError(t, err, "fake error")
		})

		t.Run("invalid status code", func(t *testing.T) {
			_, err := checkForRetry(&http.Response{StatusCode: 500}, nil)
			assert.EqualError(t, err, "HTTP 500")
		})
	})
}

func Test_retry(t *testing.T) {
	var counter int

	err := retry(func(attempt int) (bool, error) {
		counter++
		require.Equal(t, counter, attempt)
		return true, errors.New("fake error")
	})
	require.EqualError(t, err, errMaxRetriesReached.Error())
	require.Equal(t, maxRetries, counter)
}
