package solus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilterUsers(t *testing.T) {
	f := FilterUsers{}

	f.ByStatus("status")

	require.Equal(t, map[string]string{
		"filter[status]": "status",
	}, f.data)
}
