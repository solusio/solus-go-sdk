package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterIPBlocks(t *testing.T) {
	f := FilterIPBlocks{}

	f.ByName("name")

	require.Equal(t, map[string]string{
		"filter[search]": "name",
	}, f.data)
}
