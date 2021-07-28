package solus

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllowInsecure(t *testing.T) {
	c := &Client{
		HTTPClient: &http.Client{
			Transport: &http.Transport{},
		},
	}
	AllowInsecure()(c)

	require.True(t, c.HTTPClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
}

func TestSetRetryPolicy(t *testing.T) {
	c := &Client{}

	SetRetryPolicy(1, time.Second)(c)

	assert.Equal(t, 1, c.Retries)
	assert.Equal(t, time.Second, c.RetryAfter)
}

type fakeLogger struct{}

func (fakeLogger) Debugf(string, ...interface{}) {}
func (fakeLogger) Errorf(string, ...interface{}) {}

func TestWithLogger(t *testing.T) {
	c := &Client{}
	l := fakeLogger{}

	WithLogger(l)(c)

	assert.Equal(t, l, c.Logger)
}

func TestEmailAndPasswordAuthenticator_Authenticate(t *testing.T) {
	authenticator := EmailAndPasswordAuthenticator{
		Email:    "test@example.com",
		Password: "Pass80rd",
	}

	t.Run("positive", func(t *testing.T) {
		credentials := Credentials{
			AccessToken: "access token",
			TokenType:   "token type",
			ExpiresAt:   "expires at",
		}

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/auth/login", r.URL.Path)
			assertRequestBody(t, r, AuthLoginRequest{
				Email:    "test@example.com",
				Password: "Pass80rd",
			})

			writeResponse(t, w, http.StatusOK, AuthLoginResponse{
				Credentials: credentials,
			})
		}))
		defer s.Close()

		u, err := url.Parse(s.URL)
		require.NoError(t, err)

		c, err := NewClient(u, authenticator)
		require.NoError(t, err)
		require.Equal(t, credentials, c.Credentials)
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("failed to make request", func(t *testing.T) {
			_, err := NewClient(&url.URL{}, authenticator, SetRetryPolicy(0, 0))
			require.EqualError(t, err, `Post "/auth/login": unsupported protocol scheme ""`)
		})

		t.Run("invalid status", func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
			}))
			defer s.Close()

			u, err := url.Parse(s.URL)
			require.NoError(t, err)

			_, err = NewClient(u, authenticator)
			require.EqualError(t, err, "HTTP POST auth/login returns 400 status code")
		})
	})
}

func TestAPITokenAuthenticator_Authenticate(t *testing.T) {
	const token = "foo"

	authenticator := APITokenAuthenticator{
		Token: token,
	}

	c, err := NewClient(&url.URL{}, authenticator)
	require.NoError(t, err)
	require.Equal(t, Credentials{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresAt:   "",
	}, c.Credentials)
}
