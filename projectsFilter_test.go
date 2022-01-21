package solus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterProjects(t *testing.T) {
	f := FilterProjects{}

	f.ByName("name")

	require.Equal(t, map[string]string{
		"filter[search]": "name",
	}, f.data)
}
