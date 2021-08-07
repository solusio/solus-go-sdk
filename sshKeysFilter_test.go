package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterSSHKeys(t *testing.T) {
	f := FilterSSHKeys{}

	f.ByName("name")

	require.Equal(t, map[string]string{
		"filter[search]": "name",
	}, f.data)
}
