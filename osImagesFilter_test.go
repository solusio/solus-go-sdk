package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterOsImages(t *testing.T) {
	f := FilterOsImages{}

	f.ByName("name")

	require.Equal(t, map[string]string{
		"filter[search]": "name",
	}, f.data)
}
