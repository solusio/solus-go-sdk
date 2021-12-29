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

type fakeData struct {
	Foo int `json:"foo"`
}

type fakeDataResponse struct {
	Data fakeData `json:"data"`
}

type fakeDataSliceResponse struct {
	Data []fakeData `json:"data"`
}

func TestClient_create(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := fakeData{1}

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			writeResponse(t, w, http.StatusCreated, expected)
		})
		defer s.Close()

		var resp fakeDataResponse

		err := createTestClient(t, s.URL).create(context.Background(), "/foo", nil, &resp)
		require.NoError(t, err)
		assert.Equal(t, expected, resp.Data)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").create(context.Background(), string([]rune{0x02}), nil, nil)
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).create(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP POST /foo returns 400 status code")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			err := createTestClient(t, s.URL).create(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP POST /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
	})
}

func TestClient_list(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := []fakeData{{1}, {2}}

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		var resp fakeDataSliceResponse

		err := createTestClient(t, s.URL).list(context.Background(), "/foo", &resp)
		require.NoError(t, err)
		assert.ElementsMatch(t, expected, resp.Data)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").list(context.Background(), string([]rune{0x02}), nil)
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).list(context.Background(), "/foo", nil)
			assert.EqualError(t, err, "HTTP GET /foo returns 400 status code")
		})
	})
}

func TestClient_get(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := fakeData{1}

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodGet, r.Method)

			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		var resp fakeDataResponse

		err := createTestClient(t, s.URL).get(context.Background(), "/foo", &resp)
		require.NoError(t, err)
		assert.Equal(t, expected, resp.Data)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").get(context.Background(), string([]rune{0x02}), nil)
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).get(context.Background(), "/foo", nil)
			assert.EqualError(t, err, "HTTP GET /foo returns 400 status code")
		})
	})
}

func TestClient_update(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := fakeData{1}

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodPut, r.Method)

			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		var resp fakeDataResponse

		err := createTestClient(t, s.URL).update(context.Background(), "/foo", nil, &resp)
		require.NoError(t, err)
		assert.Equal(t, expected, resp.Data)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").update(context.Background(), string([]rune{0x02}), nil, nil)
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).update(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP PUT /foo returns 400 status code")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			err := createTestClient(t, s.URL).update(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP PUT /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
	})
}

func TestClient_patch(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		expected := fakeData{1}

		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodPatch, r.Method)

			writeResponse(t, w, http.StatusOK, expected)
		})
		defer s.Close()

		var resp fakeDataResponse

		err := createTestClient(t, s.URL).patch(context.Background(), "/foo", nil, &resp)
		require.NoError(t, err)
		assert.Equal(t, expected, resp.Data)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").patch(context.Background(), string([]rune{0x02}), nil, nil)
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPatch, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).patch(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP PATCH /foo returns 400 status code")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPatch, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			err := createTestClient(t, s.URL).patch(context.Background(), "/foo", nil, nil)
			assert.EqualError(t, err, "HTTP PATCH /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
	})
}

func TestClient_asyncDelete(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			writeResponse(t, w, http.StatusOK, fakeTask)
		})
		defer s.Close()

		task, err := createTestClient(t, s.URL).asyncDelete(context.Background(), "/foo")
		require.NoError(t, err)
		assert.Equal(t, fakeTask, task)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			_, err := createTestClient(t, "/").asyncDelete(context.Background(), string([]rune{0x02}))
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP DELETE /foo returns 400 status code")
		})

		t.Run("invalid JSON", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{"))
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, `failed to decode "{": unexpected end of JSON input`)
		})

		t.Run("task without id", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				writeResponse(t, w, http.StatusOK, Task{})
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, "task doesn't have an id")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP DELETE /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
	})
}

func TestClient_syncDelete(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodDelete, r.Method)

			w.WriteHeader(http.StatusNoContent)
		})
		defer s.Close()

		err := createTestClient(t, s.URL).syncDelete(context.Background(), "/foo")
		require.NoError(t, err)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			err := createTestClient(t, "/").syncDelete(context.Background(), string([]rune{0x02}))
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			err := createTestClient(t, s.URL).syncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP DELETE /foo returns 400 status code")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodDelete, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			err := createTestClient(t, s.URL).syncDelete(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP DELETE /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
	})
}

func TestClient_asyncPost(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/foo", r.URL.Path)
			assert.Equal(t, http.MethodPost, r.Method)

			writeResponse(t, w, http.StatusOK, fakeTask)
		})
		defer s.Close()

		task, err := createTestClient(t, s.URL).asyncPost(context.Background(), "/foo")
		require.NoError(t, err)
		assert.Equal(t, fakeTask, task)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			_, err := createTestClient(t, "/").asyncPost(context.Background(), string([]rune{0x02}))
			assert.EqualError(
				t,
				err,
				`failed to build HTTP request: parse "\x02": net/url: invalid control character in URL`,
			)
		})

		t.Run("invalid status code", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(http.StatusBadRequest)
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncPost(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP POST /foo returns 400 status code")
		})

		t.Run("invalid JSON", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("{"))
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncPost(context.Background(), "/foo")
			assert.EqualError(t, err, `failed to decode "{": unexpected end of JSON input`)
		})

		t.Run("task without id", func(t *testing.T) {
			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				writeResponse(t, w, http.StatusOK, Task{})
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncPost(context.Background(), "/foo")
			assert.EqualError(t, err, "task doesn't have an id")
		})

		t.Run("422 unprocessable entity", func(t *testing.T) {
			expectedErrors := map[string][]string{
				"foo": {"fizz", "buzz"},
				"bar": {"foo"},
			}

			s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/foo", r.URL.Path)
				assert.Equal(t, http.MethodPost, r.Method)

				writeJSON(t, w, http.StatusUnprocessableEntity, map[string]interface{}{
					"message": "fake error",
					"errors":  expectedErrors,
				})
			})
			defer s.Close()

			_, err := createTestClient(t, s.URL).asyncPost(context.Background(), "/foo")
			assert.EqualError(t, err, "HTTP POST /foo returns 422 status code with errors: fake error")

			var e HTTPError
			assert.ErrorAs(t, err, &e)
			assert.Equal(t, expectedErrors, e.Errors)
		})
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

		t.Run("failed to build URL", func(t *testing.T) {
			var data struct {
				Foo int
			}

			_, err := cl.buildRequest(
				context.Background(),
				http.MethodGet,
				string([]rune{0x01}),
				withBody(data),
			)
			assert.EqualError(t, err, `parse "\x01": net/url: invalid control character in URL`)
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

	t.Run("positive", func(t *testing.T) {
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
	})

	t.Run("negative", func(t *testing.T) {
		_, err := cl.buildURL(string([]rune{0x02}), requestOpts{})
		assert.EqualError(t, err, `parse "\x02": net/url: invalid control character in URL`)
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

func Test_unmarshal(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		var data struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		}

		err := unmarshal([]byte(`
{
	"foo": "fizz",
	"bar": 42
}
`), &data)
		require.NoError(t, err)
		assert.Equal(t, "fizz", data.Foo)
		assert.Equal(t, 42, data.Bar)
	})

	t.Run("negative", func(t *testing.T) {
		var data struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		}

		err := unmarshal([]byte("invalid"), &data)
		assert.EqualError(t, err, `failed to decode "invalid": invalid character 'i' looking for beginning of value`)
	})
}
