package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterUsers(t *testing.T) {
	f := FilterUsers{}

	f.ByStatus("status")

	require.Equal(t, map[string]string{
		"filter[status]": "status",
	}, f.data)
}
