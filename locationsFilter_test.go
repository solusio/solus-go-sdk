package solus

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilterLocations(t *testing.T) {
	f := FilterLocations{}

	f.ByName("name")

	require.Equal(t, map[string]string{
		"filter[search]": "name",
	}, f.data)
}
