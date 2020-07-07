package solus

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type authenticator struct{}

func (authenticator) Authenticate(*Client) (Credentials, error) { return Credentials{}, nil }

func TestAllowInsecure(t *testing.T) {
	c := &Client{
		HttpClient: &http.Client{
			Transport: &http.Transport{},
		},
	}
	AllowInsecure()(c)

	require.True(t, c.HttpClient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
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
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			_ = r.Body.Close()
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/auth/login", r.URL.Path)
			require.Equal(t, string(b), `{"email":"test@example.com","password":"Pass80rd"}`)

			b, err = json.Marshal(AuthLoginResponseData{
				Data: AuthLoginResponse{
					Credentials: credentials,
				},
			})
			require.NoError(t, err)

			w.WriteHeader(200)
			_, err = w.Write(b)
			require.NoError(t, err)
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
		require.EqualError(t, err, "HTTP 400: fake error")
	})
}
