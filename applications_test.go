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

func TestApplicationsService_Create(t *testing.T) {
	expected := Application{}
	data := ApplicationCreateRequest{
		Name:             "name",
		Url:              "url",
		IconId:           1,
		CloudInitVersion: "cloud init version",
		UserDataTemplate: "user data template",
		JsonSchema:       "json schema",
		IsVisible:        true,
		LoginLink: LoginLink{
			Type:    LoginLinkTypeURL,
			Content: "login link content",
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := ApplicationCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/applications", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, data, d)

		b, err = json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(201)
		_, _ = w.Write(b)
	})
	defer s.Close()

	u, err := url.Parse(s.URL)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)

	l, err := c.Applications.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}