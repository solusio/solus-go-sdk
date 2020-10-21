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

func TestPlansService_Create(t *testing.T) {
	expected := Plan{}
	data := PlanCreateRequest{
		Name: "name",
		Type: "type",
		Params: PlanParams{
			Disk: 1,
			RAM:  2,
			VCPU: 3,
		},
		StorageType:        "storage type",
		ImageFormat:        "image format",
		IsVisible:          true,
		IsSnapshotsEnabled: true,
		Limits: PlanLimits{
			TotalBytes: PlanLimit{
				IsEnabled: true,
				Limit:     4,
			},
			TotalIops: PlanLimit{
				IsEnabled: true,
				Limit:     5,
			},
		},
		TokenPerHour:  6,
		TokenPerMonth: 7,
		Position:      8,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)

		d := PlanCreateRequest{}
		err = json.Unmarshal(b, &d)
		require.NoError(t, err)

		require.Equal(t, "/plans", r.URL.Path)
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

	l, err := c.Plans.Create(context.Background(), data)
	require.NoError(t, err)
	require.Equal(t, expected, l)
}
