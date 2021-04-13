package solus

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("fake error"))
		}))
		defer s.Close()

		u, err := url.Parse(s.URL)
		require.NoError(t, err)

		_, err = NewClient(u, authenticator)
		require.EqualError(t, err, "HTTP POST auth/login returns 400 status code: fake error")
	})
}
