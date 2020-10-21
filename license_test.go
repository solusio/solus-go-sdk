package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestLicenseService_Activate(t *testing.T) {
	expected := License{}
	data := LicenseActivateRequest{
		ActivationCode: "activation code",
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := LicenseActivateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/license/activate", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, data, d)

		b, err = json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.License.Activate(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}
