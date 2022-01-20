package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterVirtualServers(t *testing.T) {
	f := FilterVirtualServers{}

	f.
		ByUserID(1337).
		ByComputeResourceID(42).
		ByStatus("status").
		ByVirtualizationType(VirtualizationTypeKVM)

	require.Equal(t, map[string]string{
		"filter[user_id]":             "1337",
		"filter[compute_resource_id]": "42",
		"filter[status]":              "status",
		"filter[virtualization_type]": string(VirtualizationTypeKVM),
	}, f.data)
}
