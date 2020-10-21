package solus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilterTasks(t *testing.T) {
	f := FilterTasks{}

	f.
		ByAction("action").
		ByStatus("status").
		ByComputeResourceID(42).
		ByComputeResourceVmID(1337)

	require.Equal(t, map[string]string{
		"filter[action]":                 "action",
		"filter[status]":                 "status",
		"filter[compute_resource_id]":    "42",
		"filter[compute_resource_vm_id]": "1337",
	}, f.data)
}
