package solus

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

var fakeActivityLogsEvent = ActivityLogs{
	ID:        1,
	Event:     ActivityLogsEventAdditionalIPCreated,
	UserEmail: "test@test.test",
	CreatedAt: "2021-12-20 16:27:19.602320414 +0000 UTC m=+0.000631810",
}

func TestActivityLogsService_List(t *testing.T) {
	expected := ActivityLogsResponse{
		Data: []ActivityLogs{
			fakeActivityLogsEvent,
		},
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/activity_logs", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)
		assertRequestQuery(t, r, url.Values{
			"filter[user_id]": []string{"1"},
			"filter[event]":   []string{string(ActivityLogsEventAdditionalIPCreated)},
		})

		writeJSON(t, w, http.StatusOK, expected)
	})
	defer s.Close()

	f := (&FilterActivityLogs{}).
		ByUserID(1).
		ByEvent(ActivityLogsEventAdditionalIPCreated)

	actual, err := createTestClient(t, s.URL).ActivityLogs.List(context.Background(), f)
	require.NoError(t, err)
	actual.service = nil
	require.Equal(t, expected, actual)
}
