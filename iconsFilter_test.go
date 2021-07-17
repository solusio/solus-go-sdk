package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterIcons(t *testing.T) {
	f := FilterIcons{}

	f.ByName("name")
	f.ByType(IconTypeApplication)

	require.Equal(t, map[string]string{
		"filter[search]": "name",
		"filter[type]":   string(IconTypeApplication),
	}, f.data)
}
