package solus

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func startTestServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = listener.Addr().String()
		h(w, r)
	}))

	err = s.Listener.Close()
	require.NoError(t, err)
	s.Listener = listener

	s.Start()

	return s
}

func writeJSON(t *testing.T, w http.ResponseWriter, statusCode int, r interface{}) {
	data, err := json.Marshal(r)
	require.NoError(t, err)

	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	require.NoError(t, err)
}

type authenticator struct{}

func (authenticator) Authenticate(*Client) (Credentials, error) { return Credentials{}, nil }

func createTestClient(t *testing.T, addr string) *Client {
	u, err := url.Parse(addr)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)
	return c
}
