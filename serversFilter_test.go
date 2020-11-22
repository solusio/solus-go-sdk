package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterServers(t *testing.T) {
	f := FilterServers{}

	f.
		ByUserID(1337).
		ByComputeResourceID(42).
		ByStatus("status")

	require.Equal(t, map[string]string{
		"filter[user_id]":             "1337",
		"filter[compute_resource_id]": "42",
		"filter[status]":              "status",
	}, f.data)
}
