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

func TestUsersService_List(t *testing.T) {
	expected := UsersResponse{}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/users", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t,
			url.Values{
				"filter[status]": []string{"status"},
			}.Encode(),
			r.URL.Query().Encode(),
		)

		b, err := json.Marshal(expected)
		require.NoError(t, err)

		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	defer s.Close()

	c := createTestClient(t, s.URL)

	f := (&FilterUsers{}).ByStatus("status")

	p, err := c.Users.List(context.Background(), f)
	require.NoError(t, err)
	p.service = nil
	require.Equal(t, expected, p)
}

func TestUsersService_Create(t *testing.T) {
	expected := User{}
	data := UserCreateRequest{
		Password:   "password",
		Email:      "email",
		Status:     "status",
		LanguageId: 1,
		Roles:      []int{2, 3},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := UserCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/users", r.URL.Path)
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

	l, err := c.Users.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}
